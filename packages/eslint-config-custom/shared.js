module.exports = {
  commonNamingConvention: [
    {
      "selector": "default",
      "format": ["camelCase"],
      "leadingUnderscore": "allow",
      "trailingUnderscore": "allow",
    },
    {
      "selector": "variable",
      "format": ["camelCase", "UPPER_CASE"],
      "leadingUnderscore": "allow",
      "trailingUnderscore": "allow",
    },
    {
      "selector": "typeLike",
      "format": ["PascalCase"],
    },
    { 
      "selector": "objectLiteralProperty",
      "format": null
    },
    {
      "selector": "variable",
      "modifiers": ["const", "exported"],
      "format": ["PascalCase", "UPPER_CASE"]
    },
    {
      "selector": "enumMember",
      "format": ["PascalCase"]
    }
  ]
}