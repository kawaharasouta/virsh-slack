package virsh

import (
	libvirt "libvirt.org/libvirt-go"
)

func list_help(args Args, conn *libvirt.Connect) ([]byte, error) {

	return nil, nil
}

func cmdlist(opt string) (libvirt.ConnectListAllDomainsFlags) {
	switch opt {
	case "all":
		return libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	case "active":
		return libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	}
	return 0
}

func list(args Args, conn *libvirt.Connect) (string, error) {
	if args.sub_options[0] == "help" {
		out, err := list_help(args, conn)
		return string(out), err
	}
	flags := cmdlist(args.sub_options[0])
	names := []byte{}
	doms, err := conn.ListAllDomains(flags)
	if err != nil {
		return "", err
	}
	for _, dom := range doms {
		name , _ := dom.GetName()
		names = append(names, name...)
		names = append(names, "\n"...)
	}
	return string(names), nil
}

