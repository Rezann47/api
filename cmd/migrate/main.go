// cmd/migrate — bağımsız migration CLI
//
// Kullanım:
//   go run ./cmd/migrate up
//   go run ./cmd/migrate down
//   go run ./cmd/migrate version
//   go run ./cmd/migrate steps -1
//   go run ./cmd/migrate force 4
//   go run ./cmd/migrate new add_notifications

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/Rezann47/YksKoc/internal/config"
	pkgmigrate "github.com/Rezann47/YksKoc/pkg/migrate"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cmd := os.Args[1]

	// "new" komutu runner gerektirmiyor
	if cmd == "new" {
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate new <name>")
		}
		runNew(os.Args[2])
		return
	}

	runner, err := pkgmigrate.New(&cfg.DB, "file://migrations")
	if err != nil {
		log.Fatalf("runner init: %v", err)
	}
	defer runner.Close() //nolint:errcheck

	switch cmd {
	case "up":
		if err := runner.Up(); err != nil {
			log.Fatalf("up: %v", err)
		}
		printVersion(runner)
		fmt.Println("✓ migrations applied")

	case "down":
		fmt.Print("⚠  Bu işlem TÜM migration'ları geri alır. Devam? [y/N]: ")
		var confirm string
		fmt.Scanln(&confirm) //nolint:errcheck
		if confirm != "y" && confirm != "Y" {
			fmt.Println("iptal")
			return
		}
		if err := runner.Down(); err != nil {
			log.Fatalf("down: %v", err)
		}
		fmt.Println("✓ all rolled back")

	case "steps":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate steps <n>  (negative = rollback)")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid n: %v", err)
		}
		if err := runner.Steps(n); err != nil {
			log.Fatalf("steps: %v", err)
		}
		printVersion(runner)

	case "version":
		printVersion(runner)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		if err := runner.Force(v); err != nil {
			log.Fatalf("force: %v", err)
		}
		fmt.Printf("✓ forced to version %d\n", v)

	default:
		fmt.Printf("unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printVersion(r *pkgmigrate.Runner) {
	v, dirty, err := r.Version()
	if err != nil {
		fmt.Println("version: (no migrations applied)")
		return
	}
	fmt.Printf("version: %d  dirty: %v\n", v, dirty)
}

// runNew golang-migrate CLI'ını çağırarak yeni dosya oluşturur.
// Gerekli: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
func runNew(name string) {
	c := exec.Command("migrate", "create", "-ext", "sql", "-dir", "migrations", "-seq", name)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatalf("migrate create: %v\nKurulum: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest", err)
	}
}

func printUsage() {
	fmt.Print(`Usage: migrate <command> [args]

Commands:
  up              Tüm bekleyen migration'ları uygula
  down            TÜM migration'ları geri al (tehlikeli!)
  steps <n>       n adım uygula (negatif = geri al)
  version         Mevcut versiyonu göster
  force <v>       Versiyonu zorla ayarla (dirty fix)
  new <name>      Yeni migration dosyası oluştur
`)
}
