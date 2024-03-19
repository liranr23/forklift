package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"os/exec"
	"regexp"

	"github.com/konveyor/forklift-controller/pkg/controller/plan/adapter/ova"
	"k8s.io/klog/v2"
	cnv "kubevirt.io/api/core/v1"
)

const DiskTargetSizePattern = "required size: ([0-9]+).*"

type arrayPaths []string

func (i *arrayPaths) String() string {
	return "paths"
}

func (i *arrayPaths) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var srcVolPath, dstVolPath, srcFormat, dstFormat, volumeMode, sourcePaths, mode, encodedVm string
	var paths arrayPaths
	var vm *cnv.VirtualMachine

	flag.StringVar(&srcVolPath, "src-path", "", "Source volume path")
	flag.StringVar(&dstVolPath, "dst-path", "", "Target volume path")
	flag.StringVar(&srcFormat, "src-format", "", "Format of the source volume")
	flag.StringVar(&dstFormat, "dst-format", "", "Format of the target volume")
	flag.StringVar(&volumeMode, "volume-mode", "", "Format of the target volume")
	flag.Var(&paths, "src-paths", "Source volume paths(string array)")
	flag.StringVar(&mode, "mode", "", "The mode of the convereter (convert/ova_export)")
	flag.StringVar(&encodedVm, "vm", "", "The VM(json encoded)")
	flag.Parse()

	klog.Info("mode: ", mode, "srcVolPath: ", srcVolPath, " dstVolPath: ", dstVolPath, " sourceFormat: ", srcFormat, " targetFormat: ", dstFormat, "sourcePaths: ", sourcePaths)
	var err error
	switch mode {
	case "convert":
		err = convert(srcVolPath, dstVolPath, srcFormat, dstFormat, volumeMode)

	case "ova_export":
		err = json.Unmarshal([]byte(encodedVm), vm)
		if err != nil {
			return
		}
		err = ovaExport(vm, )
	}
	if err != nil {
		klog.Fatal(err)
	}
}

func convert(srcVolPath, dstVolPath, srcFormat, dstFormat, volumeMode string) error {
	err := qemuimgConvert(srcVolPath, dstVolPath, srcFormat, dstFormat)
	if err != nil {
		return err
	}

	klog.Info("Copying over source")

	// Copy dst over src
	switch volumeMode {
	case "Block":
		err = qemuimgConvert(dstVolPath, srcVolPath, dstFormat, dstFormat)
		if err != nil {
			return err
		}
	case "Filesystem":
		// Use mv for files as it's faster than qemu-img convert
		cmd := exec.Command("mv", dstVolPath, srcVolPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr // Capture stderr
		klog.Info("Executing command: ", cmd.String())
		err := cmd.Run()
		if err != nil {
			klog.Error(stderr.String())
			return err
		}
	}

	return nil
}

func qemuimgConvert(srcVolPath, dstVolPath, srcFormat, dstFormat string) error {
	cmd := exec.Command(
		"qemu-img",
		"convert",
		"-p",
		"-f", srcFormat,
		"-O", dstFormat,
		srcVolPath,
		dstVolPath,
	)

	klog.Info("Executing command: ", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			klog.Error(line)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		klog.Info(line)
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func ovaExport(vm *cnv.VirtualMachine, srcVolPath string) (err error) {
	size, err := qemuimgMeasure(srcVolPath)
	if err != nil {
		return err
	}
	ova.WriteOvf(vm)
}

func qemuimgMeasure(srcVolPath string) (size string, err error) {
	cmd := exec.Command(
		"qemu-img",
		"measure",
		"-O", "qcow2",
		srcVolPath,
	)
	klog.Info("Executing command: ", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			klog.Error(line)
		}
	}()
	regex, err := regexp.Compile(DiskTargetSizePattern)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		size = regex.FindString(line)
		klog.Info(line)
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	return
}
