package main

import (
    "fmt"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "log"
    "regexp"
)

var (
    pcapFile string = "/vagrant/tcpdump_09222016.pcap11"
    handle   *pcap.Handle
    err      error
)

func main() {
    handle, err = pcap.OpenOffline(pcapFile)
    if err != nil { log.Fatal(err) }
    defer handle.Close()
    
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        printPacketInfo(packet)
    }
}

func printPacketInfo(packet gopacket.Packet) {

    applicationLayer := packet.ApplicationLayer()
    if applicationLayer != nil {
        re := regexp.MustCompile("Referer:(.*)\n")
        fmt.Println(re.FindStringSubmatch(string(applicationLayer.Payload())))
    }
    
    if err := packet.ErrorLayer(); err != nil {
        fmt.Println("Error decoding some part of the packet:", err)
    }
}
