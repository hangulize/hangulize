#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const _ = require('lodash')

// Read descriptions from descriptions.yaml
const descriptionsPath = path.join(__dirname, '../src/descriptions.yaml')
const yaml = require('js-yaml')
const descriptions = yaml.load(fs.readFileSync(descriptionsPath, 'utf8'))

// Read base template from public/index.html
const templatePath = path.join(__dirname, '../public/index.html')
const template = fs.readFileSync(templatePath, 'utf8')

Object.keys(descriptions).forEach((lang) => {
  const description = descriptions[lang]
  if (!description) {
    console.warn(`No data found for language: ${lang}`)
    return
  }

  // Create language-specific HTML
  const descriptionSection = `<section class="description">
${_.escape(description)}
</section>
<div class="ui divider"></div>`

  // Replace the intro section in template with description + intro
  const introSection = '<section class="intro">'
  const htmlContent = template.replace(introSection, descriptionSection + '\n' + introSection)

  // Create directory structure LANG/index.html in build directory
  const langDir = path.join(__dirname, `../build/${lang}`)
  if (!fs.existsSync(langDir)) {
    fs.mkdirSync(langDir, { recursive: true })
  }

  const outputPath = path.join(langDir, 'index.html')
  fs.writeFileSync(outputPath, htmlContent, 'utf8')
  console.log(`Generated build/${lang}/index.html`)
})

console.log(`Generated static pages for ${Object.keys(descriptions).length} languages`)
