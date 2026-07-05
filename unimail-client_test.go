package main

import "testing"

func TestIsPackageManagedInstallIncludesBrew(t *testing.T) {
	for _, method := range []string{"apt", "winget", "brew", "homebrew"} {
		if !isPackageManagedInstall(method) {
			t.Fatalf("expected %q to be package managed", method)
		}
	}
}

func TestPackageManagerInstallMethodFromPathDetectsHomebrew(t *testing.T) {
	tests := []string{
		"/opt/homebrew/Cellar/unimail-client/0.0.2/bin/unimail-client",
		"/opt/homebrew/bin/unimail-client",
		"/home/linuxbrew/.linuxbrew/Cellar/unimail-client/0.0.2/bin/unimail-client",
		"/home/linuxbrew/.linuxbrew/bin/unimail-client",
	}

	for _, path := range tests {
		if got := packageManagerInstallMethodFromPath(path); got != "brew" {
			t.Fatalf("packageManagerInstallMethodFromPath(%q) = %q, want brew", path, got)
		}
	}
}

func TestPackageManagerInstallMethodFromPathIgnoresDirectPath(t *testing.T) {
	if got := packageManagerInstallMethodFromPath("/usr/local/bin/unimail-client"); got != "" {
		t.Fatalf("packageManagerInstallMethodFromPath() = %q, want empty", got)
	}
	if got := packageManagerInstallMethodFromPath("/tmp/cellar/unimail-client"); got != "" {
		t.Fatalf("packageManagerInstallMethodFromPath() = %q, want empty", got)
	}
}
