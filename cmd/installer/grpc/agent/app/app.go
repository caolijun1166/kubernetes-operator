/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"net"
	"time"

	"github.com/gosoon/glog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	installerv1 "github.com/gosoon/kubernetes-operator/pkg/apis/installer/v1"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster"
	"github.com/gosoon/kubernetes-operator/pkg/installer/cluster/constants"
)

// TODO: this
type Flagpole struct {
	Config    string
	ImageName string
	Retain    bool
	Wait      time.Duration
	Port      string

// NewServerCommand returns a new cobra.Command for kube-on-kube server
func NewServerCommand() *cobra.Command {
	flags := &Flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cluster",
		Short: "Creates a local Kubernetes cluster",
		Long:  "Creates a local Kubernetes cluster ",
		Run: func(cmd *cobra.Command, args []string) {
			run(flags)
		},
	}
	cmd.Flags().StringVar(&flags.ImageName, "image", "", "node docker image to use for booting the cluster")
	cmd.Flags().StringVar(&flags.Registry, "registry", "registry.cn-hangzhou.aliyuncs.com/aliyun_kube_system", "kubernetes image registry")
	cmd.Flags().DurationVar(&flags.Wait, "wait", time.Duration(0), "Wait for control plane node to be ready (default 0s)")
	cmd.Flags().StringVar(&flags.Port, "port", "10023", "installer agent grpc server port(default 10023)")
	return cmd
}

func run(flags *Flagpole) {
	// start grpc server
	l, err := net.Listen("tcp", ":"+flags.Port)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	// create a cluster context and create the cluster
	ctx := cluster.NewContext(constants.DefaultClusterName, server, flags.Port)

	installer := NewInstaller(&Options{
		Flags:   flags,
		Context: ctx,
		Server:  server,
	})

	// register grpc server
	installerv1.RegisterInstallerServer(server, installer)

	glog.Fatal(server.Serve(l))
}