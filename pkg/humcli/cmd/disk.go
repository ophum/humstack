package cmd

import (
	"fmt"
	"log"

	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/spf13/cobra"
)

var diskListCmd = &cobra.Command{
	Use: "disks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newDiskClient()
		if err != nil {
			return err
		}
		disks, err := c.List(cmd.Context())
		if err != nil {
			log.Println("failed to get disks, ", err.Error())
			return err
		}
		fmt.Println("name\tsize\tstatus")
		for _, disk := range disks {
			fmt.Printf("%s\t%s\t%s\n", disk.Name, disk.LimitSize, disk.Status)
		}
		return nil
	},
}

func init() {
	diskCreateCmd.Flags().StringP("name", "n", "", "new disk name")
	diskCreateCmd.Flags().StringP("limit", "l", "", "new disk limit size")
}

var diskCreateCmd = &cobra.Command{
	Use: "disk",
	RunE: func(cmd *cobra.Command, args []string) error {

		c, err := newDiskClient()
		if err != nil {
			return err
		}

		limitSize, _ := cmd.Flags().GetString("limit")
		size, err := entity.ParseByteUnit(limitSize)
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		disk, err := c.Create(cmd.Context(), &entity.Disk{
			Name:      name,
			LimitSize: *size,
		})
		if err != nil {
			return err
		}
		fmt.Println("name\tsize\tstatus")
		fmt.Printf("%s\t%s\t%s\n", disk.Name, size, disk.Status)
		return nil
	},
}
