export default (x: number): boolean => {
  const str = x.toString();
  for (
    let i = 0, j = str.length - 1;
    i <= str.length / 2, j >= str.length / 2;
    i++, j--
  ) {
    if (str[i] !== str[j]) return false;
  }
  return true;
};
