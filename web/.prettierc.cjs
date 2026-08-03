module.exports = {
  plugins: [require.resolve("@trivago/prettier-plugin-sort-imports")],
  printWidth: 80,
  trailingComma: "all",
  tabWidth: 2,
  semi: true,
  singleQuote: false,
  importOrder: ["^[^src].*$", "^src/(.*)[^.css]$", "^(.*).css$"],
  importOrderSeparation: true,
  importOrderSortSpecifiers: true,
};
