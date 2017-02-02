package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Masterminds/glide/action"
	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/msg"
	gpath "github.com/Masterminds/glide/path"
)

func main() {
	var (
		glideYaml  = gpath.DefaultGlideFile
		argDebug   = false
		argVerbose = false
		argQuiet   = false
		argExact   = false
	)

	flag.StringVar(&glideYaml, "yaml", gpath.DefaultGlideFile, "Set a YAML configuration file")
	flag.BoolVar(&argVerbose, "verbose", false, "Print more verbose informational messages")
	flag.BoolVar(&argDebug, "debug", false, "Print debug verbose informational messages")
	flag.BoolVar(&argQuiet, "quiet", false, "Quiet (no info or debug messages)")
	flag.BoolVar(&argExact, "exact", false, "Exact (replace version ranges with specific versions)")
	flag.Parse()

	action.Verbose(argVerbose)
	action.Debug(argDebug)
	action.Quiet(argQuiet)
	gpath.GlideFile = glideYaml

	base := "."

	msg.Verbose("Loading Glide config from %s...", glideYaml)
	conf := action.EnsureConfig()

	if !gpath.HasLock(base) {
		msg.Die("Lock file (glide.lock) does not exist. Please run update.")
		return
	}

	lock, err := cfg.ReadLockFile(filepath.Join(base, gpath.LockFile))
	if err != nil {
		msg.Die("Could not load lockfile.")
	}

	hash, err := conf.Hash()
	if err != nil {
		msg.Die("Could not load lockfile.")
	} else if hash != lock.Hash {
		msg.Warn("Lock file may be out of date. Hash check of YAML failed. You may need to run 'update'")
	}

	msg.Verbose("Checking for unpinned packages...")

	newconf := conf.Clone()
	newconf.Imports = dependenciesFromLocks(lock.Imports)
	newconf.DevImports = dependenciesFromLocks(lock.DevImports)

	if !argExact {
		restoreYAMLVersions(newconf.Imports, conf.Imports)
		restoreYAMLVersions(newconf.DevImports, conf.DevImports)
	}

	var newhash string
	newhash, err = newconf.Hash()
	if err != nil {
		msg.Die("Could not marshal config.")
	} else if newhash == hash {
		msg.Info("Everything is already pinned!")
		os.Exit(0)
	}

	msg.Verbose("Pinning locked packages and subpackages...")

	glideYamlFile, _ := gpath.Glide()
	if err = newconf.WriteFile(glideYamlFile); err != nil {
		msg.Die("Error while write Glide config to file back: %v", err)
	}

	msg.Info("New Glide config has been updated with pinned packages.")
}

func loadGlideConfig() (config *cfg.Config) {
	glideYamlFile, err := gpath.Glide()
	if err != nil {
		msg.Die("Could not find Glide config file")
	}

	var yml []byte
	if yml, err = ioutil.ReadFile(glideYamlFile); err != nil {
		msg.Die("Error while reading config file: %v", err)
	}

	if config, err = cfg.ConfigFromYaml(yml); err != nil {
		msg.Die("Error while parsing config file: %v", err)
	}

	return config
}

func dependenciesFromLocks(locks cfg.Locks) cfg.Dependencies {
	deps := make(cfg.Dependencies, len(locks))
	for i, l := range locks {
		dep := cfg.DependencyFromLock(l)
		deps[i] = dep
	}

	return deps
}

func restoreYAMLVersions(newdeps, olddeps cfg.Dependencies) {
	for _, dep := range newdeps {
		if olddep := olddeps.Get(dep.Name); olddep != nil {
			if olddep.Reference != "" {
				dep.Reference = olddep.Reference
			}
		}
	}
}
