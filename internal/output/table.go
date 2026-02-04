package output

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/tehnerd/vape/internal/models"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(data interface{}, writer io.Writer) error {
	switch v := data.(type) {
	case []models.Measurement:
		return f.formatMeasurements(v, writer)
	case *models.Measurement:
		return f.formatMeasurement(v, writer)
	case []models.Probe:
		return f.formatProbes(v, writer)
	case *models.Probe:
		return f.formatProbe(v, writer)
	case *models.Credits:
		return f.formatCredits(v, writer)
	case []models.MeasurementResult:
		return f.formatResults(v, writer)
	default:
		return fmt.Errorf("unsupported type for table formatting: %T", data)
	}
}

func (f *TableFormatter) formatMeasurements(measurements []models.Measurement, writer io.Writer) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"ID", "Type", "Target", "Status", "Probes", "Created"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, m := range measurements {
		created := time.Unix(m.CreationTime, 0).Format("2006-01-02 15:04")
		table.Append([]string{
			strconv.Itoa(m.ID),
			string(m.Type),
			m.Target,
			string(m.Status),
			strconv.Itoa(m.ParticipantCount),
			created,
		})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatMeasurement(m *models.Measurement, writer io.Writer) error {
	created := time.Unix(m.CreationTime, 0).Format("2006-01-02 15:04:05")
	started := ""
	if m.StartTime > 0 {
		started = time.Unix(m.StartTime, 0).Format("2006-01-02 15:04:05")
	}
	stopped := ""
	if m.StopTime > 0 {
		stopped = time.Unix(m.StopTime, 0).Format("2006-01-02 15:04:05")
	}

	fmt.Fprintf(writer, "ID:           %d\n", m.ID)
	fmt.Fprintf(writer, "Type:         %s\n", m.Type)
	fmt.Fprintf(writer, "Target:       %s\n", m.Target)
	fmt.Fprintf(writer, "Status:       %s\n", m.Status)
	fmt.Fprintf(writer, "AF:           %d\n", m.AF)
	fmt.Fprintf(writer, "Description:  %s\n", m.Description)
	fmt.Fprintf(writer, "One-off:      %t\n", m.IsOneoff)
	fmt.Fprintf(writer, "Probes:       %d requested, %d scheduled, %d participating\n",
		m.ProbesRequested, m.ProbesScheduled, m.ParticipantCount)
	fmt.Fprintf(writer, "Created:      %s\n", created)
	if started != "" {
		fmt.Fprintf(writer, "Started:      %s\n", started)
	}
	if stopped != "" {
		fmt.Fprintf(writer, "Stopped:      %s\n", stopped)
	}
	if len(m.ResolvedIPs) > 0 {
		fmt.Fprintf(writer, "Resolved IPs: %v\n", m.ResolvedIPs)
	}

	return nil
}

func (f *TableFormatter) formatProbes(probes []models.Probe, writer io.Writer) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"ID", "Country", "ASN", "Status", "IPv4", "IPv6"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, p := range probes {
		asnStr := ""
		if p.ASNv4 > 0 {
			asnStr = strconv.Itoa(p.ASNv4)
		}
		ipv4 := p.AddressV4
		if ipv4 == "" {
			ipv4 = "-"
		}
		ipv6 := p.AddressV6
		if ipv6 == "" {
			ipv6 = "-"
		}
		if len(ipv6) > 20 {
			ipv6 = ipv6[:17] + "..."
		}
		table.Append([]string{
			strconv.Itoa(p.ID),
			p.CountryCode,
			asnStr,
			string(p.Status.Name),
			ipv4,
			ipv6,
		})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatProbe(p *models.Probe, writer io.Writer) error {
	fmt.Fprintf(writer, "ID:          %d\n", p.ID)
	fmt.Fprintf(writer, "Country:     %s\n", p.CountryCode)
	fmt.Fprintf(writer, "Status:      %s\n", p.Status.Name)
	fmt.Fprintf(writer, "Is Anchor:   %t\n", p.IsAnchor)
	fmt.Fprintf(writer, "Is Public:   %t\n", p.IsPublic)
	if p.Description != "" {
		fmt.Fprintf(writer, "Description: %s\n", p.Description)
	}
	fmt.Fprintf(writer, "\nIPv4:\n")
	fmt.Fprintf(writer, "  Address:   %s\n", p.AddressV4)
	fmt.Fprintf(writer, "  Prefix:    %s\n", p.PrefixV4)
	fmt.Fprintf(writer, "  ASN:       %d\n", p.ASNv4)
	if p.AddressV6 != "" {
		fmt.Fprintf(writer, "\nIPv6:\n")
		fmt.Fprintf(writer, "  Address:   %s\n", p.AddressV6)
		fmt.Fprintf(writer, "  Prefix:    %s\n", p.PrefixV6)
		fmt.Fprintf(writer, "  ASN:       %d\n", p.ASNv6)
	}
	if p.Latitude != 0 || p.Longitude != 0 {
		fmt.Fprintf(writer, "\nLocation:    %.4f, %.4f\n", p.Latitude, p.Longitude)
	}
	if len(p.Tags) > 0 {
		fmt.Fprintf(writer, "\nTags:\n")
		for _, t := range p.Tags {
			fmt.Fprintf(writer, "  - %s\n", t.Name)
		}
	}

	return nil
}

func (f *TableFormatter) formatCredits(c *models.Credits, writer io.Writer) error {
	fmt.Fprintf(writer, "Current Balance:         %d\n", c.CurrentBalance)
	fmt.Fprintf(writer, "Estimated Daily Income:  %d\n", c.EstimatedDailyIncome)
	fmt.Fprintf(writer, "Estimated Daily Expense: %d\n", c.EstimatedDailyExpense)
	fmt.Fprintf(writer, "Estimated Daily Balance: %d\n", c.EstimatedDailyBalance)
	if c.EstimatedRunoutSeconds > 0 {
		days := c.EstimatedRunoutSeconds / 86400
		fmt.Fprintf(writer, "Estimated Runout:        %d days\n", days)
	}
	return nil
}

func (f *TableFormatter) formatResults(results []models.MeasurementResult, writer io.Writer) error {
	if len(results) == 0 {
		fmt.Fprintln(writer, "No results available")
		return nil
	}

	// Determine result type from first result
	switch results[0].Type {
	case "ping":
		return f.formatPingResults(results, writer)
	case "traceroute":
		return f.formatTracerouteResults(results, writer)
	case "dns":
		return f.formatDNSResults(results, writer)
	default:
		return f.formatGenericResults(results, writer)
	}
}

func (f *TableFormatter) formatPingResults(results []models.MeasurementResult, writer io.Writer) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Probe", "From", "Target", "Sent", "Rcvd", "Min", "Avg", "Max"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, r := range results {
		table.Append([]string{
			strconv.Itoa(r.ProbeID),
			r.From,
			r.DstAddr,
			strconv.Itoa(r.Sent),
			strconv.Itoa(r.Rcvd),
			fmt.Sprintf("%.2f", r.Min),
			fmt.Sprintf("%.2f", r.Avg),
			fmt.Sprintf("%.2f", r.Max),
		})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatTracerouteResults(results []models.MeasurementResult, writer io.Writer) error {
	for i, r := range results {
		if i > 0 {
			fmt.Fprintln(writer, "---")
		}
		fmt.Fprintf(writer, "Probe %d (%s) -> %s\n", r.ProbeID, r.From, r.DstAddr)
		for _, hop := range r.Hops {
			if hop.Error != "" {
				fmt.Fprintf(writer, "  %2d: error: %s\n", hop.Hop, hop.Error)
				continue
			}
			for _, result := range hop.Result {
				if result.Timeout {
					fmt.Fprintf(writer, "  %2d: *\n", hop.Hop)
				} else if result.Error != "" {
					fmt.Fprintf(writer, "  %2d: %s (error: %s)\n", hop.Hop, result.From, result.Error)
				} else {
					fmt.Fprintf(writer, "  %2d: %s  %.3f ms\n", hop.Hop, result.From, result.RTT)
				}
			}
		}
	}
	return nil
}

func (f *TableFormatter) formatDNSResults(results []models.MeasurementResult, writer io.Writer) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Probe", "From", "Resolver", "RT (ms)", "Answers"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, r := range results {
		if r.DNSResult == nil {
			continue
		}
		answerCount := len(r.DNSResult.Answers)
		table.Append([]string{
			strconv.Itoa(r.ProbeID),
			r.From,
			r.DNSResult.DstAddr,
			fmt.Sprintf("%.2f", r.DNSResult.Time),
			strconv.Itoa(answerCount),
		})
	}

	table.Render()
	return nil
}

func (f *TableFormatter) formatGenericResults(results []models.MeasurementResult, writer io.Writer) error {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Probe", "From", "Target", "Type", "Timestamp"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, r := range results {
		ts := time.Unix(r.Timestamp, 0).Format("2006-01-02 15:04:05")
		table.Append([]string{
			strconv.Itoa(r.ProbeID),
			r.From,
			r.DstAddr,
			r.Type,
			ts,
		})
	}

	table.Render()
	return nil
}
