package os

import (
	"fmt"
	"strings"

	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type YumPackageManager struct{}

func (p YumPackageManager) Install(pkg ...string) string {
	return fmt.Sprintf("yum install -y %s", strings.Join(pkg, " "))
}

func (p YumPackageManager) Update() string {
	return ""
}
func (p YumPackageManager) Uninstall(pkg ...string) string {
	return fmt.Sprintf("yum remove -y %s", strings.Join(pkg, " "))
}
func (p YumPackageManager) Mark(pkg ...string) string {
	return ""
}
func (p YumPackageManager) ListInstalled() string {
	return ""
}
func (p YumPackageManager) CleanupCaches() string {
	return ""
}

func (p YumPackageManager) GetInstalledVersion(pkg string) string {
	pkg = strings.Split(pkg, "=")[0]
	stdout, ok := utils.SafeExec("yum info " + pkg)
	if !ok {
		log.Debugf("Failed installation check for %s -> %s", pkg, stdout)
		return ""
	}

	if !strings.Contains(stdout, "Installed Packages") {
		log.Debugf("%s is in db, but is not installed", pkg)
		return ""
	}

	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "Version") {
			return strings.Split(line, ": ")[1]
		}
	}
	log.Debugf("Unable to find version info in " + stdout)
	return ""
}

func (p YumPackageManager) AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string) string {
	repo := fmt.Sprintf(
		`[%s]
name=%s
baseurl=%s
enabled=1
`, name, name, url)

	if gpgKey != "" {
		repo += fmt.Sprintf(`gpgcheck=1
repo_gpgcheck=1
gpgkey=%s`, gpgKey)
	}
	return fmt.Sprintf(`cat <<EOF >/etc/yum.repos.d/%s.repo
%s
EOF`, name, repo)
}

func (p YumPackageManager) Setup() string {
	return "[[ $(which curl 2> /dev/null) && $(which sudo 2> /dev/null) ]] || yum install -y sudo curl"
}
