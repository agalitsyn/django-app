package main

import (
	"fmt"
	"net"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/rigging"
	"github.com/gravitational/trace"
)

func getRPCClient(stolonRPCHost, stolonRPCPort string) (*rpc.Client, error) {
	rpcEndpoint := net.JoinHostPort(stolonRPCHost, stolonRPCPort)
	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return client, nil
}

func install(stolonRPCHost, stolonRPCPort, dbName string) error {
	log.Infof("Creating database '%s'", dbName)
	client, err := getRPCClient(stolonRPCHost, stolonRPCPort)
	if err != nil {
		return trace.Wrap(err, "Dialing to stolon's RPC failed")
	}

	var reply string
	command := "DatabaseOperation.Create"
	err = client.Call(command, dbName, &reply)
	log.Infof("Execute RPC command '%s' on stolon's RPC", command)
	if err != nil {
		return trace.Wrap(err, fmt.Sprintf("Can't create database '%s'", dbName))
	}
	log.Infof("Reply: %s", reply)

	log.Infof("Creating django service and replication controller")
	out, err := rigging.FromFile(rigging.ActionCreate, "/var/lib/gravity/resources/django.yaml")
	if err != nil {
		return trace.Wrap(err, fmt.Sprintf("cmd output: %s", string(out)))
	}

	return nil
}

func uninstall(stolonRPCHost, stolonRPCPort, dbName string) error {
	log.Infof("Deleting django service and replication controller")
	out, err := rigging.FromFile(
		rigging.ActionDelete,
		"/var/lib/gravity/resources/django.yaml")
	if err != nil {
		return trace.Wrap(err, fmt.Sprintf("cmd output: %s", string(out)))
	}

	log.Infof("Deleting database '%s'", dbName)
	client, err := getRPCClient(stolonRPCHost, stolonRPCPort)
	if err != nil {
		return trace.Wrap(err, "Dialing to stolon's RPC failed")
	}

	var reply string
	command := "DatabaseOperation.Delete"
	err = client.Call(command, dbName, &reply)
	log.Infof("Execute RPC command '%s' on stolon's RPC", command)
	if err != nil {
		return trace.Wrap(err, fmt.Sprintf("Can't delete database '%s'", dbName))
	}

	log.Infof("Reply: %s", reply)
	return nil
}
