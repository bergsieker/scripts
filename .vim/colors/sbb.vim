" Vim color file
" Maintainer Steven Bergsieker

hi clear
set background=dark
if exists("syntax_on")
  syntax reset
endif

let g:colors_name = "sbb"

hi Normal cterm=None ctermbg=Black ctermfg=Gray

" General appearance
hi CursorLine cterm=None ctermfg=DarkMagenta
hi DiffAdd ctermbg=DarkBlue ctermfg=Gray
hi DiffChange ctermbg=DarkBlue ctermfg=Gray
hi DiffDelete ctermbg=DarkBlue ctermfg=DarkGray
hi DiffText ctermbg=DarkBlue ctermfg=DarkGreen
hi IncSearch ctermbg=DarkBlue ctermfg=DarkGreen
hi LineNr ctermfg=DarkCyan
hi Search ctermbg=DarkBlue ctermfg=DarkGreen
hi StatusLine cterm=Bold ctermbg=DarkCyan ctermfg=Black
hi StatusLineNC cterm=Underline ctermbg=Black ctermfg=DarkCyan
hi VertSplit cterm=Bold ctermbg=DarkCyan ctermfg=Black

" Syntax highlighting
hi Comment ctermfg=DarkCyan
hi Constant ctermfg=DarkRed
hi String ctermfg=DarkGreen
hi PreProc ctermfg=DarkYellow
hi PreCondit ctermfg=DarkRed
hi Statement ctermfg=DarkYellow
hi Type ctermfg=DarkYellow
