package cmd

import (
	"errors"
	"path/filepath"

	"github.com/xalanq/cf-tool/client"
	"github.com/xalanq/cf-tool/config"
)

// Parse command
func Parse() (err error) {
	cfg := config.Instance
	cln := client.Instance
	info := Args.Info
	source := ""
	ext := ""
	script := ""
	if cfg.GenAfterParse {
		if len(cfg.Template) == 0 {
			return errors.New("You have to add at least one code template by `cf config`")
		}
		path := cfg.Template[cfg.Default].Path
		ext = filepath.Ext(path)
		if source, err = readTemplateSource(path, cln); err != nil {
			return
		}
		if script, err = readTemplateSource(cfg.ScriptTemplatePath, cln); err != nil {
			return
		}
	}
	work := func() error {
		_, paths, err := cln.Parse(info)
		if err != nil {
			return err
		}
		if cfg.GenAfterParse {
			for _, path := range paths {
				gen(source, path, ext)
				if err := genScript(script, path, ext); err != nil {
					println(err)
				}
			}
		}
		return nil
	}
	if err = work(); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = work()
		}
	}
	return
}
