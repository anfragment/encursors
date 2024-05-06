export function getFlagEmoji(countryCode) {
  const codePoints = countryCode.toUpperCase().split('').map(char =>  0x1F1E6 + char.charCodeAt(0) - 'A'.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}
