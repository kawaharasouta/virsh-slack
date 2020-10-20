package virsh

import (
	"log"

	libvirt "libvirt.org/libvirt-go"
)

var commands = map[string](func(Args, *libvirt.Connect) (string, error)){
	"help":     help,
	"list":     list,
	"start":    start,
	"shutdown": shutdown,
}

type Args struct {
	sub_command string
	sub_options []string
}

func help(args Args, conn *libvirt.Connect) (string, error) {
	return "help", nil
}

func start(args Args, conn *libvirt.Connect) (string, error) {
	return "start", nil
}

func shutdown(args Args, conn *libvirt.Connect) (string, error) {
	return "shutdown", nil
}

func Exec_Virsh(args Args) (string, error){
	if command, ok := commands[args.sub_command]; ok {
		conn, err := libvirt.NewConnect("qemu:///system")
		if err != nil {
			log.Println(err)
			return "", err
		}
		defer conn.Close()
		out, err := command(args, conn)
		return out, err
	}
	return "", nil
}

func virsh_Args(command []string, args *Args) {
	// if (strings.HasPrefix(command[0], "-")) {
	//  if (strings.HasPrefix(command[0][1:], "-")) {
	//    fmt.Println("long options")
	//    args.options = append(args.options, command[0][2:])
	//    fmt.Println(args)
	//  } else {
	//    fmt.Println("short options")
	//  }
	// } else {
	//  fmt.Println("subcommand")
	// }
	args.sub_command = command[0]
	args.sub_options = command[1:]
}

func Virsh(command []string) (string, error) {
	var args Args
	virsh_Args(command, &args)
	out, err := Exec_Virsh(args)
	return out, err
}



