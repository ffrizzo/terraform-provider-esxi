package esxi

import (
	"fmt"
	"log"
  "strconv"
  "strings"
)


func resourcePoolCREATE(c *Config, resource_pool_name string, cpu_min int,
  cpu_min_expandable bool, cpu_max int, cpu_shares string, mem_min int,
  mem_min_expandable bool, mem_max int, mem_shares string, parent_pool string) (string, error) {

  esxiSSHinfo := SshConnectionStruct{c.Esxi_hostname, c.Esxi_hostport, c.Esxi_username, c.Esxi_password}
	log.Println("[provider-esxi / resourcePoolCREATE]" )
	var pool_id, remote_cmd string

  cpu_min_opt := ""
  if cpu_min > 0 {
    cpu_min_opt = fmt.Sprintf("--cpu-min=%d", cpu_min)
  }

  cpu_min_expandable_opt := "--cpu-min-expandable=true"
  if cpu_min_expandable == false {
    cpu_min_expandable_opt = "--cpu-min-expandable=false"
  }

  cpu_max_opt := ""
  if cpu_max > 0 {
    cpu_max_opt = fmt.Sprintf("--cpu-max=%d", cpu_max)
  }

  cpu_shares_opt := "--cpu-shares=normal"
  if cpu_shares == "low" ||  cpu_shares == "high" {
    cpu_shares_opt = fmt.Sprintf("--cpu-shares=%s", cpu_shares)
  } else {
    tmp_var, err := strconv.Atoi(cpu_shares)
    if err == nil {
      cpu_shares_opt = fmt.Sprintf("--cpu-shares=%d", tmp_var)
    }
  }

  mem_min_opt := ""
  if mem_min > 0 {
    mem_min_opt = fmt.Sprintf("--mem-min=%d", mem_min)
  }

  mem_min_expandable_opt := "--mem-min-expandable=true"
  if mem_min_expandable == false {
    mem_min_expandable_opt = "--mem-min-expandable=false"
  }

  mem_max_opt := ""
  if mem_max > 0 {
    mem_max_opt = fmt.Sprintf("--mem-max=%d", mem_max)
  }

  mem_shares_opt := "--mem-shares=normal"
  if mem_shares == "low" ||  mem_shares == "high" {
    mem_shares_opt = fmt.Sprintf("--mem-shares=%s", mem_shares)
  } else {
    tmp_var, err := strconv.Atoi(mem_shares)
    if err == nil {
      mem_shares_opt = fmt.Sprintf("--mem-shares=%d", tmp_var)
    }
  }

  parent_pool_id, err := getPoolID(c, parent_pool)
  if err != nil {
    return "", err
  }

  remote_cmd = fmt.Sprintf("vim-cmd hostsvc/rsrc/create %s %s %s %s %s %s %s %s %s %s",
    cpu_min_opt,cpu_min_expandable_opt, cpu_max_opt, cpu_shares_opt,
    mem_min_opt,mem_min_expandable_opt, mem_max_opt, mem_shares_opt, parent_pool_id, resource_pool_name)

	stdout,_ := runRemoteSshCommand(esxiSSHinfo, remote_cmd, "create resource pool")
  log.Printf("[provider-esxi / resourcePoolCREATE] stdout |%s|", stdout)
  pool_id, err = getPoolID(c, resource_pool_name)
  if err == nil {
    return pool_id, err
  }

  r := strings.NewReplacer("'vim.ResourcePool:","", "'","")
  stdout = r.Replace(stdout)
  stdout = strings.TrimSpace(stdout)
  return stdout, err
}