# asn1
encoding/asn1 has bugs on bigint encoding, according to  https://www.itu.int/ITU-T/studygroups/com17/languages/X.690-0207.pdf, Section 8.3.2.
so we rewrite the asn1 from go 1.15.7 sources.
