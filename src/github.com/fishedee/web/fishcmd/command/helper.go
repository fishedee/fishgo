package command

import (
	"fmt"
	"github.com/fishedee/web/fishcmd/modules"
	"strings"
	"time"
)

var (
	stdPackage = map[string]bool{
		"archive":    true,
		"bufio":      true,
		"builtin":    true,
		"bytes":      true,
		"compress":   true,
		"container":  true,
		"crypto":     true,
		"database":   true,
		"debug":      true,
		"encoding":   true,
		"errors":     true,
		"expvar":     true,
		"flag":       true,
		"fmt":        true,
		"go":         true,
		"hash":       true,
		"html":       true,
		"image":      true,
		"index":      true,
		"io":         true,
		"log":        true,
		"math":       true,
		"mime":       true,
		"net":        true,
		"os":         true,
		"path":       true,
		"reflect":    true,
		"regexp":     true,
		"runtime":    true,
		"sort":       true,
		"strconv":    true,
		"strings":    true,
		"sync":       true,
		"syscall":    true,
		"testing":    true,
		"text":       true,
		"time":       true,
		"unicode":    true,
		"unsafe":     true,
		"gopkg.in":   true,
		"golang.org": true,
	}
)

func buildAllNormal(appName string) error {
	//安装文件
	var err error
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + appName + ")...")
	timer.Start()
	err = modules.InstallPackage(appName)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}

	//复制文件
	err = modules.CopyFile(modules.GetAppInstallPath(), modules.GetAppCurrentPath())
	if err != nil {
		modules.Log.Error("copy fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

func checkPackageCode(packageName string, codeModifyTime time.Time) (bool, error) {
	packageLibraryInfo, err := modules.GetPackageLibraryInfo(packageName)
	if err != nil {
		return true, nil
	}
	if packageLibraryInfo.ModifyTime.After(codeModifyTime) {
		return false, nil
	} else {
		return true, nil
	}
}

func isStdPackage(packageName string) bool {
	firstIndex := strings.Index(packageName, "/")
	if firstIndex != -1 {
		packageName = packageName[0:firstIndex]
	}
	_, isExist := stdPackage[packageName]
	return isExist
}

func buildSinglePackage(packageName string) (bool, error) {
	var err error
	oldPackageLibaryInfo, err := modules.GetPackageLibraryInfo(packageName)
	if err != nil {
		oldPackageLibaryInfo = modules.PackageLibraryInfo{}
	}

	err = modules.InstallPackage(packageName)
	if err != nil {
		return false, err
	}

	newPackageLibraryInfo, err := modules.GetPackageLibraryInfo(packageName)
	if err != nil {
		return false, err
	}

	if oldPackageLibaryInfo.Symbol == newPackageLibraryInfo.Symbol {
		return false, nil
	} else {
		return true, nil
	}

}

func buildPackages(packageName string, buildingPackage map[string]bool, buildedPackage map[string]bool) (bool, error) {
	if isStdPackage(packageName) == true {
		return false, nil
	}
	if _, isExist := buildingPackage[packageName]; isExist == true {
		return false, fmt.Errorf("import cycle %v", packageName)
	}
	if _, isExist := buildedPackage[packageName]; isExist == true {
		return buildedPackage[packageName], nil
	}
	buildingPackage[packageName] = true

	packageCodeInfo, err := modules.GetPackageCodeInfo(packageName)
	if err != nil {
		return false, err
	}

	dependenceIsChange := false
	for _, singleDependence := range packageCodeInfo.Dependence {
		singleIsChange, err := buildPackages(singleDependence, buildingPackage, buildedPackage)
		if err != nil {
			return false, err
		}
		if singleIsChange {
			dependenceIsChange = true
		}
	}

	if packageCodeInfo.Name != "main" {
		//非入口包
		codeIsChange, err := checkPackageCode(packageName, packageCodeInfo.ModifyTime)
		if err != nil {
			return false, err
		}
		delete(buildingPackage, packageName)
		if dependenceIsChange || codeIsChange {
			isChange, err := buildSinglePackage(packageName)
			if err != nil {
				return false, err
			}
			buildedPackage[packageName] = isChange
			return isChange, nil
		} else {
			err := modules.RefreshPackageLibrary(packageName)
			if err != nil {
				return false, err
			}
			buildedPackage[packageName] = false
			return false, nil
		}
	} else {
		//入口包
		err = modules.InstallPackage(packageName)
		if err != nil {
			return false, err
		}

		err = modules.CopyFile(modules.GetAppInstallPath(), modules.GetAppCurrentPath())
		if err != nil {
			return false, err
		}

		return false, nil
	}
}

func buildAllFast(appName string) error {
	timer := modules.NewTimer()
	modules.Log.Debug("start building (" + appName + ")...")
	timer.Start()
	buildingPackage := map[string]bool{}
	buildedPackage := map[string]bool{}
	_, err := buildPackages(appName, buildingPackage, buildedPackage)
	if err != nil {
		modules.Log.Error("build fail! error: %v", err.Error())
		return err
	}
	timer.Stop()
	modules.Log.Debug("build success! time: %v", timer.Elapsed())
	return nil
}

type builder func(packageName string) error

func getBuildAll(argv []string) builder {
	if len(argv) >= 1 && argv[0] == "-fast" {
		return buildAllFast
	} else {
		return buildAllNormal
	}
}

func run(appName string) error {
	err := modules.RunPackage(appName)
	if err != nil {
		modules.Log.Error("%v running fail! error: %v", appName, err.Error())
		return err
	}
	modules.Log.Debug("%v is running", appName)
	return nil
}
