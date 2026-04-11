package main

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/Viswesh934/gotei/internal/engine"
)

type EPFNominee struct {
	NameAndAddress string
	Relationship   string
	DateOfBirth    string
	SharePercent   string
	Guardian       string
}

type EPSFamilyMember struct {
	Name         string
	Address      string
	DateOfBirth  string
	Relationship string
}

type EPSNominee struct {
	NameAndAddress string
	DateOfBirth    string
	Relationship   string
	SharePercent   string
}

type Form2Data struct {
	Name                string
	FatherOrHusbandName string
	DateOfBirth         string
	Sex                 string
	MaritalStatus       string
	AccountNumber       string
	PermanentAddress    string
	TemporaryAddress    string
	DateOfJoiningEPF    string
	DateOfJoiningEPS    string
	EmployerPlace       string
	EmployerDate        string
	LogoPath            string
	EmployerSignature   string
	EmployerSeal        string
	AdobeSig            string
	EPFNominees         []EPFNominee
	EPSFamilyMembers    []EPSFamilyMember
	EPSNominees         []EPSNominee
}

type viewEPSNominee struct {
	NameAndAddress string
	DateOfBirth    string
	Relationship   string
}

type form2View struct {
	Form2Data
	EPFNominees      []EPFNominee
	EPSFamilyMembers []EPSFamilyMember
	EPSNominees      []viewEPSNominee
	Today            string
}

func main() {
	outDir := "output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(err)
	}

	logoPath := filepath.Join(outDir, "image.png")
	signPath := filepath.Join(outDir, "employer-signature.png")
	sealPath := filepath.Join(outDir, "employer-seal.png")

	data := sampleData(logoPath, signPath, sealPath)
	html, err := renderForm2HTML(data)
	if err != nil {
		panic(err)
	}

	pdf, err := engine.Render(html)
	if err != nil {
		panic(err)
	}

	outPDF := filepath.Join(outDir, "form2-revised-fullscale.pdf")
	if err := os.WriteFile(outPDF, pdf, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("ok %s\n", outPDF)
}

func sampleData(logoPath, signPath, sealPath string) Form2Data {
	return Form2Data{
		Name:                "Aarav Kumar",
		FatherOrHusbandName: "Ramesh Kumar",
		DateOfBirth:         "14/08/1993",
		Sex:                 "Male",
		MaritalStatus:       "Married",
		AccountNumber:       "MH/BAN/0001122/0000789",
		PermanentAddress:    "22 Park Avenue, Lucknow, Uttar Pradesh",
		TemporaryAddress:    "Flat 9C, Green Residency, Pune, Maharashtra",
		DateOfJoiningEPF:    "03/06/2021",
		DateOfJoiningEPS:    "03/06/2021",
		EmployerPlace:       "Pune",
		EmployerDate:        time.Now().Format("02/01/2006"),
		LogoPath:            logoPath,
		EmployerSignature:   signPath,
		EmployerSeal:        sealPath,
		AdobeSig:            "{{BigSig_es_:signer1:signature:dimension(width=35mm,height=6mm)}}",
		EPFNominees: []EPFNominee{
			{NameAndAddress: "Meera Kumar, Flat 9C, Pune", Relationship: "Spouse", DateOfBirth: "01/01/1995", SharePercent: "70", Guardian: "-"},
			{NameAndAddress: "Arjun Kumar, Flat 9C, Pune", Relationship: "Son", DateOfBirth: "12/03/2020", SharePercent: "30", Guardian: "Meera Kumar (Mother)"},
		},
		EPSFamilyMembers: []EPSFamilyMember{
			{Name: "Meera Kumar", Address: "Flat 9C, Pune", DateOfBirth: "01/01/1995", Relationship: "Spouse"},
			{Name: "Arjun Kumar", Address: "Flat 9C, Pune", DateOfBirth: "12/03/2020", Relationship: "Son"},
		},
		EPSNominees: []EPSNominee{
			{NameAndAddress: "Meera Kumar, Flat 9C, Pune", DateOfBirth: "01/01/1995", Relationship: "Spouse", SharePercent: "100"},
		},
	}
}

func renderForm2HTML(data Form2Data) (string, error) {
	view := form2View{
		Form2Data:        data,
		EPFNominees:      padEPF(data.EPFNominees, 3),
		EPSFamilyMembers: padEPSFamily(data.EPSFamilyMembers, 3),
		EPSNominees:      padEPSNominees(data.EPSNominees, 3),
		Today:            time.Now().Format("02/01/2006"),
	}

	tpl := htmltemplate.Must(htmltemplate.New("form2").Funcs(htmltemplate.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}).Parse(form2Template))

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func padEPF(in []EPFNominee, n int) []EPFNominee {
	if len(in) > n {
		in = in[:n]
	}
	out := append([]EPFNominee{}, in...)
	for len(out) < n {
		out = append(out, EPFNominee{})
	}
	return out
}

func padEPSFamily(in []EPSFamilyMember, n int) []EPSFamilyMember {
	if len(in) > n {
		in = in[:n]
	}
	out := append([]EPSFamilyMember{}, in...)
	for len(out) < n {
		out = append(out, EPSFamilyMember{})
	}
	return out
}

func padEPSNominees(in []EPSNominee, n int) []viewEPSNominee {
	if len(in) > n {
		in = in[:n]
	}
	out := make([]viewEPSNominee, 0, n)
	for _, item := range in {
		rel := item.Relationship
		if item.SharePercent != "" {
			rel = fmt.Sprintf("%s (%s%%)", rel, item.SharePercent)
		}
		out = append(out, viewEPSNominee{
			NameAndAddress: item.NameAndAddress,
			DateOfBirth:    item.DateOfBirth,
			Relationship:   rel,
		})
	}
	for len(out) < n {
		out = append(out, viewEPSNominee{})
	}
	return out
}

const form2Template = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Form 2 Revised - Nomination & Declaration</title>
</head>
<body style="font-family: Arial, sans-serif; font-size: 12px; line-height: 1.35; color: #000; margin: 0; padding: 12px;">
<table style="width: 100%; border-collapse: collapse; margin-bottom: 8px;">
	<tr>
		<td style="width: 86px; vertical-align: top;">
			<img src="{{.LogoPath}}" style="width: 48px; height: 48px; display: block; margin: 0;" />
		</td>
		<td style="vertical-align:top;" align="center">
			<div style="font-weight: 700; font-size: 14px;">FORM - 2 ( Revised )</div>
			<div style="font-weight: 700; font-size: 13px; margin-top: 2px;">NOMINATION AND DECLARATION FORM</div>
			<div style="font-size: 11px; margin-top: 2px;">FOR EXEMPTED / UNEXEMPTED ESTABLISHMENTS</div>
			<div style="font-size: 10px; margin-top: 2px;">Declaration and Nomination Form Under the Employee's Provident Funds &amp; Employees' Pension Scheme</div>
			<div style="font-size: 9px; margin-top: 2px;">Paragraph 33 &amp; 61(1) of the Employees' Provident Fund Scheme, 1952 and Paragraph 18 of the Employees' Pension Scheme, 1995</div>
		</td>
	</tr>
</table>

<div style="width: 100%; border: 0.5px solid #000; padding: 8px 8px 4px; margin-top: 6px; box-sizing: border-box;">
<table style="width: 100%; border-collapse: collapse; border-spacing: 0 4px; font-size: 12px;" cellspacing="0">
	<tr><td style="width: 34px; padding: 2px 4px 2px 0; vertical-align: top;">1</td><td style="padding: 2px 8px; vertical-align: top;">Name ( In Block Letters )</td><td style="width: 14px; padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.Name}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">2</td><td style="padding: 2px 8px; vertical-align: top;">Father's / Husband's Name</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.FatherOrHusbandName}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">3</td><td style="padding: 2px 8px; vertical-align: top;">Date of Birth</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.DateOfBirth}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">4</td><td style="padding: 2px 8px; vertical-align: top;">Sex</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.Sex}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">5</td><td style="padding: 2px 8px; vertical-align: top;">Marital Status</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.MaritalStatus}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">6</td><td style="padding: 2px 8px; vertical-align: top;">Account Number</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.AccountNumber}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">7</td><td style="padding: 2px 8px; vertical-align: top;">Address</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">7a</td><td style="padding: 2px 8px; vertical-align: top;">Permanent</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.PermanentAddress}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">7b</td><td style="padding: 2px 8px; vertical-align: top;">Temporary</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.TemporaryAddress}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">8a</td><td style="padding: 2px 8px; vertical-align: top;">Date of Joining EPF</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.DateOfJoiningEPF}}</span></td></tr>
	<tr><td style="padding: 2px 4px 2px 0; vertical-align: top;">8b</td><td style="padding: 2px 8px; vertical-align: top;">Date of Joining EPS</td><td style="padding: 2px 4px; vertical-align: top; text-align: center;">:</td><td style="padding: 2px 0 2px 8px; vertical-align: top;"><span style="color: #0B66FF;">{{.DateOfJoiningEPS}}</span></td></tr>
</table>
</div>

<div style="text-align: center; font-weight: 700; margin-top: 10px;">PART - A ( EPF )</div>
<p style="font-size: 10px; margin: 6px 0 0 0; width: 100%; text-align: center;">I hereby nominate the person(s) / cancel the nomination made by me previously and person(s) mentioned below to receive the amount standing to my credit in the Employees' Provident Fund, in the event of my death.</p>
<table border="0" cellspacing="0" cellpadding="4" style="width:100%; border-collapse:collapse; font-size: 10px;">
<tr>
	<th style="border:0.5px solid #000;">S.NO</th>
	<th style="border:0.5px solid #000;">Name &amp; Address of the Nominee / Nominees</th>
	<th style="border:0.5px solid #000;">Relationship</th>
	<th style="border:0.5px solid #000;">Date of Birth</th>
	<th style="border:0.5px solid #000;">Share of accumulation</th>
	<th style="border:0.5px solid #000;">Guardian details (if minor)</th>
</tr>
{{range $i, $n := .EPFNominees}}
<tr>
	<td style="border:0.5px solid #000;">{{add1 $i}}</td>
	<td style="height:34px; border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.NameAndAddress}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.Relationship}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.DateOfBirth}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.SharePercent}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.Guardian}}</span></td>
</tr>
{{end}}
</table>
<p style="width:100%; text-align: left;">Certified that I have no family as defined in para 2(g) of the Employees' Provident Fund Scheme, 1952 and should I acquire a family hereafter this nomination shall stand cancelled.Certified that my father / mother is / are dependent upon me.Unmarried members in the absence of dependent parents may nominate any other person to receive the amount.</p>
<div style="margin-top:8px; font-size: 10px;">Signature / Thumb impression of the Subscriber: ________________________________ <span style="color:#0B66FF;">{{.AdobeSig}}</span></div>

<div style="text-align: center; font-weight: 700; margin-top: 10px;">PART - B ( EPS )</div>
<p style="font-size: 10px; margin: 6px 0 0 0; width: 100%; text-align: left;">I hereby furnish below particulars of the members of my family who would be eligible to receive widow / children pension in the event of my death.</p>
<table border="0" cellspacing="0" cellpadding="4" style="width:100%; border-collapse:collapse; font-size: 10px;">
<tr><th style="border:0.5px solid #000;">S.No</th><th style="border:0.5px solid #000;">Name</th><th style="border:0.5px solid #000;">Address</th><th style="border:0.5px solid #000;">Date of Birth</th><th style="border:0.5px solid #000;">Relationship</th></tr>
{{range $i, $m := .EPSFamilyMembers}}
<tr>
	<td style="border:0.5px solid #000;">{{add1 $i}}</td>
	<td style="height:34px; border:0.5px solid #000;"><span style="color: #0B66FF;">{{$m.Name}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$m.Address}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$m.DateOfBirth}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$m.Relationship}}</span></td>
</tr>
{{end}}
</table>


<p style="margin: 0 0 6px 0; width: 100%; text-align: left;">I hereby nominate the following person for receiving monthly widow pension in the event of my death without leaving any eligible family member for receiving pension.</p>
<table border="0" cellspacing="0" cellpadding="4" style="width:100%; border-collapse:collapse; font-size: 10px;">
<tr><th style="border:0.5px solid #000;">Name &amp; Address</th><th style="border:0.5px solid #000;">Date of Birth</th><th style="border:0.5px solid #000;">Relationship</th></tr>
{{range $i, $n := .EPSNominees}}
<tr>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.NameAndAddress}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.DateOfBirth}}</span></td>
	<td style="border:0.5px solid #000;"><span style="color: #0B66FF;">{{$n.Relationship}}</span></td>
</tr>
{{end}}
</table>

<div style="margin-top:8px; font-size: 10px;">Date: <span style="color: #0B66FF;">{{.Today}}</span></div>
<div style="margin-top:6px; font-size: 10px;">Signature / Thumb impression of the subscriber: ________________________________ <span style="color:#0B66FF;">{{.AdobeSig}}</span></div>

<div style="text-align: center; font-weight: 700; margin-top: 10px;">CERTIFICATE BY EMPLOYER</div>
<p style="font-size: 10px; margin: 6px 0 0 0; width: 100%; text-align: left;">Certified that the above declaration and nomination has been signed/thumb impressed by the employee after the entries have been read over and explained.</p>
<div style="margin-top:8px; font-size: 10px;">Place: <span style="color: #0B66FF;">{{.EmployerPlace}}</span></div>
<div style="margin-top:4px; font-size: 10px;">Date: <span style="color: #0B66FF;">{{.EmployerDate}}</span></div>
<div style="margin-top:6px; font-size: 10px;">Signature of Employer: ________________________________</div>
<div style="margin-top:6px; font-size: 10px;">Name &amp; Address of the Establishment: ________________________________</div>
<table style="margin-top:10px; border-collapse:collapse;">
	<tr>
		<td style="padding-right:16px;"><img src="{{.EmployerSignature}}" style="width:100px; height:40px; border:1px solid #333;" /></td>
		<td><img src="{{.EmployerSeal}}" style="width:80px; height:80px; border:1px solid #333;" /></td>
	</tr>
</table>
</body>
</html>`
