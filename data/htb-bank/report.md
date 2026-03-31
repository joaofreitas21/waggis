---
TAGS: TODO (comma-separated: e.g., SMB, RCE)
DESCRIPTION: TODO
---

# hack the box legacy 

**Platform:** hackthebox  
**Difficulty:** easy  
**Date:** 2026-03-18

---

## Summary


---

## Reconnaissance

### Port Scanning

Command used:
```bash
nmap -sC -sV -oN nmap/initial 
```

Results:

![Nmap scan results](nmap.png)

Key findings:


### Service Enumeration

Further enumeration of discovered services.

```bash
# Commands used
```

![Service enumeration](enum.png)

---

## Initial Access

### Vulnerability Identification

Description of the vulnerability found.

### Exploitation

Steps taken to exploit the vulnerability:

```bash
# Exploit commands
```

![Exploitation proof](exploit.png)

User flag: `USER_FLAG_HERE`

---

## Privilege Escalation

### Enumeration

Commands run for privilege escalation enumeration:

```bash
# Enumeration commands
```

![Privilege escalation enumeration](privesc-enum.png)

### Exploitation

Description of privilege escalation method:

```bash
# Privesc commands
```

![Root access](root.png)

Root flag: `ROOT_FLAG_HERE`

---

## Tools Used


---

## References

- [Exploit DB link](#)
- [CVE details](#)
- [Other resources](#)