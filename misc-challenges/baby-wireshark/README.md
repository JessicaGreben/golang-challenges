# Baby Wireshark

### Desscription:

We have recorded a packet capture of an HTTP request and response for an image, performed over an imperfect network. 

The challenge is to parse the capture file, located at `net.cap`, find and parse the packets constituting the image download, and reconstruct the image!

### Details

`tcpdump` was used to capture the network traffic and create the `net.cap` file.

The file is saved as “pcap-savefile” format. Read more about that format here:
https://www.tcpdump.org/manpages/pcap-savefile.5.txt

