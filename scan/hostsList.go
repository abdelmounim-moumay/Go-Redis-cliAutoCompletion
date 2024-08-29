package scan

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("Host already in the list")
	ErrNotExists = errors.New("host no in the list")
)

type HostsList struct {
	Hosts []string
}

// search hosts
func (hl *HostsList) search(host string) (bool, int) {

	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1

}

// add hosts

func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// remove hosts
func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// load obtains hosts
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return nil
}

// save hosts
func (hl *HostsList) Save(hostsFile string) error {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	return ioutil.WriteFile(hostsFile, []byte(output), 0644)
}
