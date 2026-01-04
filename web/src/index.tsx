import 'semantic-ui-less/semantic.less'

import * as React from 'react'
import * as ReactDOM from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router-dom'

import App from './components/App'

const rootElem = document.getElementById('root') as HTMLElement
const introHTML = rootElem.getElementsByClassName('intro')[0].innerHTML
const root = ReactDOM.createRoot(rootElem)
root.render(
  <BrowserRouter>
    <Routes>
      <Route path="/:lang" element={<App introHTML={introHTML} />} />
      <Route path="/:lang/:word" element={<App introHTML={introHTML} />} />
      <Route path="/" element={<App introHTML={introHTML} />} />
    </Routes>
  </BrowserRouter>
)
