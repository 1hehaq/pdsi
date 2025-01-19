<div align="center">

<h3>
  <b>
    
  🪶 <kbd>[**PDSI**](https://github.com/1hehaq/pdsi)</kbd>
  
  </b>
</h3>

<h6>Read sensitive information in PDF documents with ease</h6>

</div>

<br>
<br>
<br>

> [!Important]
> **_PDSI may show you false positives, because it doesn't have human brain..!! So always put your manual effort to verify the data._**

<br>
<br>
<br>

|                                      |                                 |
| :----------------------------------: | :-----------------------------: |
|              **`direct input`**            |             **`stdin`**             |
| ![direct input](https://github.com/user-attachments/assets/25f4d6b3-380e-47ce-a893-0051b058fa36) |  ![STDin](https://github.com/user-attachments/assets/d98872b1-fddd-4a54-b4b6-99ba913e0fd3)   |

<br>
<br>

<h6 align="center">Install</h6>

```bash
go install github.com/1hehaq/pdsi@latest
```

<br>
<br>

<h6 align="center">
  Inbuilt keywords
</h6>

<pre>
  <i>
confidential, private, restricted, internal, not for distribution, do not share, proprietary, trade secret, classified
sensitive, bank statement, invoice, salary, contract, agreement, non disclosure, passport, social security, ssn, date of birth
credit card, identity, id number, internal use only, company confidential, staff only, management only, internal only
  </i>
</pre>

<br>
<br>
<br>

<h6 align="center">
  example commands
</h6>


```bash
cat pdf.txt | pdsi
```

<div align="center">
<kbd>OR</kbd>
</div>

<br>

```bash
pdsi -pdf pdf.txt
```

<div align="center">
<kbd>OR</kbd>
</div>

<br>

```bash
pdsi -local file.pdf -match "secret,confidential,internal"
```

<div align="center">
<kbd>OR</kbd>
</div>

<br>

```bash
cat list-of-pdf.txt | pdsi
```

<div align="center">
<kbd>OR</kbd>
</div>

<br>

```bash
pdsi -local document.pdf # "document-1.pdf,document-2.pdf,document-3.pdf"
```


<br>
<br>
<br>

<h6 align="center">kindly for hackers</h6>


<div align="center">
  <a href="https://github.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/github.png" alt="GitHub"></a>
  <a href="https://twitter.com/1hehaq"><img src="https://img.icons8.com/material-outlined/20/808080/twitter.png" alt="X"></a>
</div>
