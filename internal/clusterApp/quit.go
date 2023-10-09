package clusterapp

import "fmt"

func (c *Cluster) Quit(port int) error {
	for i, target := range c.Targets {
		if target.srv.Addr == fmt.Sprintf("localhost:%d", port) {
			c.Targets = append(c.Targets[0:i], c.Targets[i+1:c.N-1]...)
			c.N--
			return nil
		}
	}

	return fmt.Errorf("error quitting target by given port")
}
