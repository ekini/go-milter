/* Author: Eugene Dementiev <eugene@dementiev.eu>
 */
package main

const (
	SMFIC_OPTNEG  string = "O"
	SMFIC_MACRO   string = "D"
	SMFIC_BODYEOB string = "E"
	SMFIC_QUIT    string = "Q"
)
const (
	SMFIF_ADDHDRS    uint32 = 0x01 << iota /* filter may add headers */
	SMFIF_CHGBODY                          /* filter may replace body */
	SMFIF_ADDRCPT                          /* filter may add recipients */
	SMFIF_DELRCPT                          /* filter may delete recipients */
	SMFIF_CHGHDRS                          /* filter may change/delete headers */
	SMFIF_QUARANTINE                       /* filter may quarantine envelope */
	/* Introduced with Sendmail 8.14. */
	SMFIF_CHGFROM     /* filter may replace sender */
	SMFIF_ADDRCPT_PAR /* filter may add recipients + args */
	SMFIF_SETSYMLIST  /* filter may send macro names */
)
const (
	SMFIP_NOCONNECT uint32 = 0x01 << iota /* filter does not want connect info */
	SMFIP_NOHELO                          /* filter does not want HELO info */
	SMFIP_NOMAIL                          /* filter does not want MAIL info */
	SMFIP_NORCPT                          /* filter does not want RCPT info */
	SMFIP_NOBODY                          /* filter does not want body */
	SMFIP_NOHDRS                          /* filter does not want headers */
	SMFIP_NOEOH                           /* filter does not want EOH */
	SMFIP_NR_HDR                          /* filter won't reply for header */
	SMFIP_NOUNKNOWN                       /* filter does not want unknown cmd */
	SMFIP_NODATA                          /* filter does not want DATA */
	/* Introduced with Sendmail 8.14. */
	SMFIP_SKIP        /* MTA supports SMFIR_SKIP */
	SMFIP_RCPT_REJ    /* filter wants rejected RCPTs */
	SMFIP_NR_CONN     /* filter won't reply for connect */
	SMFIP_NR_HELO     /* filter won't reply for HELO */
	SMFIP_NR_MAIL     /* filter won't reply for MAIL */
	SMFIP_NR_RCPT     /* filter won't reply for RCPT */
	SMFIP_NR_DATA     /* filter won't reply for DATA */
	SMFIP_NR_UNKN     /* filter won't reply for UNKNOWN */
	SMFIP_NR_EOH      /* filter won't reply for eoh */
	SMFIP_NR_BODY     /* filter won't reply for body chunk */
	SMFIP_HDR_LEADSPC /* header value has leading space */
)
const (
	SMFIP_NOSEND_MASK uint32 = SMFIP_NOCONNECT | SMFIP_NOHELO | SMFIP_NOMAIL | SMFIP_NORCPT | SMFIP_NOBODY | SMFIP_NOHDRS | SMFIP_NOEOH | SMFIP_NOUNKNOWN | SMFIP_NODATA

	SMFIP_NOREPLY_MASK uint32 = SMFIP_NR_CONN | SMFIP_NR_HELO | SMFIP_NR_MAIL | SMFIP_NR_RCPT | SMFIP_NR_DATA | SMFIP_NR_UNKN | SMFIP_NR_HDR | SMFIP_NR_EOH | SMFIP_NR_BODY
)

type SMFIC_OPTIONS struct {
	Version  uint32
	Actions  uint32
	Protocol uint32
}
