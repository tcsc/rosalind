package codon

type codon [3]byte

func New() codon {
	return codon{'\x00', '\x00', '\x00'}
}

var Table = map[codon]rune{
	codon{'U', 'U', 'U'}: 'F',
	codon{'C', 'U', 'U'}: 'L',
	codon{'A', 'U', 'U'}: 'I',
	codon{'G', 'U', 'U'}: 'V',

	codon{'U', 'U', 'C'}: 'F',
	codon{'C', 'U', 'C'}: 'L',
	codon{'A', 'U', 'C'}: 'I',
	codon{'G', 'U', 'C'}: 'V',

	codon{'U', 'U', 'A'}: 'L',
	codon{'C', 'U', 'A'}: 'L',
	codon{'A', 'U', 'A'}: 'I',
	codon{'G', 'U', 'A'}: 'V',

	codon{'U', 'U', 'G'}: 'L',
	codon{'C', 'U', 'G'}: 'L',
	codon{'A', 'U', 'G'}: 'M',
	codon{'G', 'U', 'G'}: 'V',

	codon{'U', 'C', 'U'}: 'S',
	codon{'C', 'C', 'U'}: 'P',
	codon{'A', 'C', 'U'}: 'T',
	codon{'G', 'C', 'U'}: 'A',

	codon{'U', 'C', 'C'}: 'S',
	codon{'C', 'C', 'C'}: 'P',
	codon{'A', 'C', 'C'}: 'T',
	codon{'G', 'C', 'C'}: 'A',

	codon{'U', 'C', 'A'}: 'S',
	codon{'C', 'C', 'A'}: 'P',
	codon{'A', 'C', 'A'}: 'T',
	codon{'G', 'C', 'A'}: 'A',

	codon{'U', 'C', 'G'}: 'S',
	codon{'C', 'C', 'G'}: 'P',
	codon{'A', 'C', 'G'}: 'T',
	codon{'G', 'C', 'G'}: 'A',

	codon{'U', 'A', 'U'}: 'Y',
	codon{'C', 'A', 'U'}: 'H',
	codon{'A', 'A', 'U'}: 'N',
	codon{'G', 'A', 'U'}: 'D',

	codon{'U', 'A', 'C'}: 'Y',
	codon{'C', 'A', 'C'}: 'H',
	codon{'A', 'A', 'C'}: 'N',
	codon{'G', 'A', 'C'}: 'D',

	codon{'U', 'A', 'A'}: '\x00',
	codon{'C', 'A', 'A'}: 'Q',
	codon{'A', 'A', 'A'}: 'K',
	codon{'G', 'A', 'A'}: 'E',

	codon{'U', 'A', 'G'}: '\x00',
	codon{'C', 'A', 'G'}: 'Q',
	codon{'A', 'A', 'G'}: 'K',
	codon{'G', 'A', 'G'}: 'E',

	codon{'U', 'G', 'U'}: 'C',
	codon{'C', 'G', 'U'}: 'R',
	codon{'A', 'G', 'U'}: 'S',
	codon{'G', 'G', 'U'}: 'G',

	codon{'U', 'G', 'C'}: 'C',
	codon{'C', 'G', 'C'}: 'R',
	codon{'A', 'G', 'C'}: 'S',
	codon{'G', 'G', 'C'}: 'G',

	codon{'U', 'G', 'A'}: '\x00',
	codon{'C', 'G', 'A'}: 'R',
	codon{'A', 'G', 'A'}: 'R',
	codon{'G', 'G', 'A'}: 'G',

	codon{'U', 'G', 'G'}: 'W',
	codon{'C', 'G', 'G'}: 'R',
	codon{'A', 'G', 'G'}: 'R',
	codon{'G', 'G', 'G'}: 'G',
}
