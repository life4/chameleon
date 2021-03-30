package chameleon

import (
	"fmt"

	"github.com/errata-ai/vale/check"
	"github.com/errata-ai/vale/core"
	"github.com/errata-ai/vale/lint"
)

type Linter struct {
	Article *Article
}

func (linter Linter) Alerts() ([]core.Alert, error) {
	config := core.NewConfig()
	config.GBaseStyles = []string{"proselint", "write-good", "Joblint", "Spelling"}
	config.MinAlertLevel = 1
	config.InExt = ".html"
	vale := lint.Linter{Config: config, CheckManager: check.NewManager(config)}
	html, err := linter.Article.HTML()
	if err != nil {
		return nil, fmt.Errorf("cannot read html: %v", err)
	}
	files, err := vale.LintString(html)
	if err != nil {
		return nil, fmt.Errorf("cannot lint html: %v", err)
	}
	return files[0].SortedAlerts(), nil
}
