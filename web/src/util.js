function getSpec(specs, lang) {
  for (let i = 0; i < specs.length; i++) {
    const spec = specs[i]
    if (spec.lang.id === lang) {
      return spec
    }
  }
  return null
}

export { getSpec }
