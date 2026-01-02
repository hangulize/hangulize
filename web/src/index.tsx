import 'semantic-ui-less/semantic.less'

import * as React from 'react'
import * as ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import App from './components/App'

const rootElem = document.getElementById('root') as HTMLElement
const introHTML = rootElem.innerHTML
const root = ReactDOM.createRoot(rootElem)
root.render(
  <BrowserRouter>
    <App introHTML={introHTML}></App>
  </BrowserRouter>
)
