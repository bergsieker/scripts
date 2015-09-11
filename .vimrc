" ##########################################################################
" ~/.vimrc
" ##########################################################################

" ==========================================================================
" Abbreviations
" ==========================================================================

:ab line- -----------------------------------------------------------------------------
:ab line= =============================================================================
:ab line/ /////////////////////////////////////////////////////////////////////////////
:ab line# #############################################################################
:ab zs_cr Copyright 2010 ZeroSoft Inc. All Rights Reserved.<CR>Author: sbberg@synopsys.com (Steven Bergsieker)
:ab snps_cr SYNOPSYS CONFIDENTIAL -- This is an unpublished, proprietary work<CR>of Synopsys, Inc., and is fully protected under copyright and trade<CR>secret laws. You may not view, use, disclose, copy, or distribute<CR>this file or any information contained herein except pursuant to a<CR>valid written license from Synopsys.
:ab FBC FIXME_BEFORE_COMMITTING
:ab TODOs TODO(sbergsieker):

" ==========================================================================
" Settings
" ==========================================================================

" Incremental search (uses only / and not :/)
set incsearch
set hlsearch
" Case insensitive unless the search string has caps in it
set smartcase
" Show matching brackets
set showmatch
" Case insensitive
set ic
" Writes a file before moving on
set autowrite
" Set syntax highlight
syntax on
" Keep 3 lines above and below cursor
set scrolloff=3
" Status line
set laststatus=2
set statusline=%<%f\ %h%m%r%=%-14.(%l,%c%V%)\ %P
" Fix backspace
set bs=2
fixdel
" Always expandtab
set expandtab
" Dictionary
set dictionary+=/usr/share/dict/words;
set dictionary+=./tags;
set complete-=k complete+=k
" For the menu with CTRL-D
" First tab: longest match, list in the statusbar.
" Next tabs: cycle through matches. (Like in the shell)
set wildmenu wildmode=longest:full,full
" Do not redraw while running macros (much faster) (LazyRedraw)
set lz
" Sharing windows clipboard
set clipboard+=unnamed
" Use numbers.
set number
" Number of tabs
set tabpagemax=40
"set textwidth=79
set textwidth=0
set switchbuf="usetab"
filetype on
filetype indent on
" Global options
set tabstop=2
set smarttab
set shiftwidth=2
set expandtab
" Show the line where the cursor is
set cursorline
" Put all the swap files in a single directory.
"let &dir="$HOME/vim_swp"

" ==========================================================================
" Shortcuts.
" ==========================================================================
" Set comment color
" hi Comment ctermfg=Cyan guifg=#80a0ff gui=bold
" Set color for diffing
" hi DiffAdd ctermfg=DarkGreen ctermbg=Blue
" hi DiffChange ctermfg=White ctermbg=Blue
" hi DiffDelete cterm=bold ctermfg=DarkMagenta ctermbg=Blue
" hi DiffText cterm=underline ctermfg=DarkRed ctermbg=Blue
" hi Pmenu ctermbg=Black
" Set color for folding
" hi Folded ctermfg=DarkMagenta ctermbg=0
" hi FoldColumn ctermfg=DarkMagenta ctermbg=0
" Set color for perl
" hi perlVarPlain ctermfg=White
" hi perlVarPlain2 ctermfg=White
"
" hi TabLineSel term=reverse cterm=bold ctermfg=7 ctermbg=1
" hi StatusLine term=reverse cterm=bold ctermfg=7 ctermbg=1
" hi ModeMsg ctermfg=DarkGreen ctermbg=Blue

" Cursor highlight.
" hi CursorLine cterm=NONE ctermfg=White ctermbg=Blue
" hi CursorLine cterm=bold ctermfg=green ctermbg=black

colorscheme sbb

set noremap
set noremap
noremap <CR> O
noremap <C-p> :tabn<CR>
noremap <C-o> :tabp<CR>
"map <C-h> 2h
"map <C-j> 2j
"map <C-k> 2k
"map <C-l> 2l
"map! <C-h> 2ha
"map! <C-j> 2ja
"map! <C-k> 2ka
"map! <C-l> 2la

" map Ctrl-A, Ctrl-E, and Ctrl-K in *all* modes.
" map! makes the mapping work in insert and commandline modes too.
map  <C-a> <Home>
map  <C-e> <End>
map! <C-a> <Home>
map! <C-e> <End>

map gd :tabe <cfile><CR>

" ==========================================================================
" Functions
" ==========================================================================

function! StartFolding()
  set foldmethod=syntax
  set foldlevel=0
  set foldenable
endfunction
:command! StartFolding call StartFolding()
"map <F7> <ESC>:StartFolding<CR>

function! StopFolding()
  set nofoldenable
endfunction
:command! StopFolding call StopFolding()
"map <S-F7> <ESC>:StopFolding<CR>

" Open a tag in a new tab.
function! TagTab()
  let currentTag=expand("<cword>")
  tabnew
  exe "tjump ".currentTag
endfunction
:command! TagTab call TagTab()
map <C-Y><C-T> <Esc>:TagTab<CR>

" Open a tag in a preview window.
function! TagPreview()
  let currentTag=expand("<cword>")
  exe "ptjump ".currentTag
endfunction
:command! TagPreview call TagPreview()
map <C-Y><C-P> <Esc>:TagPreview<CR>

" Open a new tab and does MRU
function! TabMRU()
  tabnew
  MRU
endfunction
:command! Tmru call TabMRU()

" Open a new tab and look for files in the same directory of this file
function! Tabe()
  tabe %
  Explore
endfunction
:command! Te call Tabe()

" Open a new tab and look for files in the same directory of this file
function! TabE()
  tabnew
  NERDTree
  on
endfunction
:command! TE call TabE()

function! SlowTerm()
  set noshowcmd
  set noruler
  set noshowmatch
  set scrolljump=5
endfunction
:command! SlowTerm call SlowTerm()

"function! Exp()
"  execute "normal I//\<Esc>"
"endfunction
":command! Exp call Exp()

function! KillSpaces()
  %s/\s\s*$//g
endfunction
:command! KS call KillSpaces()

function! KillTabs()
  %s/	/        /g
endfunction
:command! KT call KillTabs()

function! KillExtra()
  /\%>81c
endfunction
:command! KE call KillExtra()

function! FixAll()
  " Fix assignment.
  "%s/\(".\{-0,}\S\)\(=\| \{2,}=\)/\1 =/cg
  "%s/\(".\{-0,}\)\(=\|= \{2,}\)\(\S\)/\1= \3/cg
  " Replace
  " "appe= hello"
  " "appe = hello"
  " "appe =hello"
  " "appe  =hello"
  " "appe  = hello"
  " "appe  =  hello" "appe= hello"
  " "appe  =  hello" "appe =hello"
  "
  "%s/"\(.*\)\(\a\+\) *= */"\1\2 = /cg
  %s/\(\s*\/\/.\{-0,}\)\<ppe\>\|\<Ppe\>/\1PPE/cg
  %s/\(\s*\/\/.\{-0,}\)\<spe\>\|\<Spe\>/\1SPE/cg
  %s/\(\s*\/\/.\{-0,}\)\<dma\>\|\<Dma\>/\1DMA/cg
  %s/\(\s*\/\/.\{-0,}\)\<dmas\>\|\<Dmas\>/\1DMAs/cg
  " Replace
  " // ppe, PPE, Ppe
  " // Ppe
  " // strunz, Ppe Ppe
  " // spe, PPE, Ppe
endfunction
:command! FA call FixAll()

function! Split() range
  "execute("'<,'>s/, /, <CR>/g")
  let n = a:firstline
  let count = 0
  while n <= a:lastline
    let line = getline(".")
    let repl = substitute(line, ', ', ", <CR>", "g")
    call setline(".", repl)
    let n = n + 1
  endwhile
endfunction
:command! Split call Split()

function! FullScreen()
  " Number of columns to display
  "let &columns=156
  set columns=158
  " Number of lines to display
  set lines=57
endfunction
:command! FS call FullScreen()

function! FullScreen2()
  " Number of columns to display
  "let &columns=156
  set columns=161
  " Number of lines to display
  set lines=58
endfunction
:command! FS2 call FullScreen2()

function! ReloadCscope()
  :cs kill -1
  ":cs reset
  " Open a cscope connection
  :cs add ./cscope.out
endfunction
:command! ReloadCscope call ReloadCscope()

function! Tags()
  :!tags.sh
  :ReloadCscope
endif
endfunction
:command! Tags call Tags()

" ==========================================================================
" Plug-ins
" ==========================================================================

" Commentify
" source ~/.vim/plugin/FeralToggleCommentify.vim

" MRC plug-in
"let MRU_File="$HOME/.vim_files_mru"
"let MRU_File="$HOME/.vim_files_mru"
" let MRU_Auto_Close=1
" let MRU_Max_Entries=500

" ctags
"set tags=./tags;

" TagList
" :help taglist.txt
" let Tlist_Show_One_File=1

" showmarks
" let showmarks_enable=0
" let showmarks_hlline_other=1
" hi ShowMarksHLl cterm=underline ctermfg=Red ctermbg=Black
" hi ShowMarksHLu cterm=underline ctermfg=Red ctermbg=Black
" hi ShowMarksHLo cterm=underline ctermfg=Red ctermbg=Black
" hi ShowMarksHLm cterm=underline ctermfg=Red ctermbg=Black
"let g:showmarks_include="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
" let g:showmarks_include="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
" hi CursorColumn cterm=bold ctermfg=White ctermbg=Black

" Order of files in Explorer windows
let g:netrw_sort_sequence="[\/]$,.am,.in,*,.swp"

" MiniBufferExplorer
"source ~/.vim/plugin/minibufexpl.vim
"map <C-o> :MBEbp<CR>
"map <C-p> :MBEbn<CR>
"map <Leader>c :CMiniBufExplorer<cr>
"map <Leader>u :UMiniBufExplorer<cr>
"map <Leader>t :TMiniBufExplorer<cr>

" NERDTree
" noremap <C-l> :NERDTreeToggle<CR>
" map! <C-l> :NERDTreeToggle<CR>

" Omnicppcomplete
"set nocp
"filetype plugin on

" ==========================================================================
" Code-style
" ==========================================================================

"autocmd FileType c,cpp :set cindent
"autocmd FileType c,cpp set omnifunc=ccomplete#CompleteCpp
"set cinoptions=>sesf0{0:sg0t0+0(0,Ws

function! SetZeroSoftBufferOptions()
  " default settings for new buffers
  setlocal autoindent
  setlocal cinoptions=g1,h1,(0
  setlocal nocindent
  setlocal shiftwidth=2
  " This used to set smartindent, but doesn't any more.
  " smartindent is evil.  filetype-detected indentexpr
  " (/usr/share/vim/vim*/indent/*.vim) is good.

  if &filetype == "make"
    setlocal shiftwidth=8

  elseif &filetype == "conf"
    " This is the filetype used for Perforce temporary files, among
    " other things
    setlocal shiftwidth=8

  elseif &filetype == 'cpp' || &filetype == 'hpp' || &filetype == 'c' || &filetype == 'h' || &filetype == 'java'
    " TODO: c-lineup-math?
    " TODO: cout << foo\n<< bar; (second << should line up with first)
    " TODO: (statement-case-open . +) ; case w/ {
    setlocal cindent
    setlocal indentexpr=ZeroSoftCIndent()
    setlocal expandtab

  elseif &filetype == 'am'
    setlocal shiftwidth=8
    setlocal noexpandtab
    setlocal list

  elseif &filetype == 'python' || &filetype == 'perl' || &filetype == 'tcl' || &filetype == 'vimrc'
    setlocal expandtab

  elseif &filetype == 'html' || &filetype == 'xml' || &filetype == 'gxp'
    setlocal expandtab

  elseif &filetype == 'txt' || &filetype == 'log'
    setlocal textwidth=79

  else

  endif

  if exists("*SetMyBufferOptions")
    call SetMyBufferOptions()
  endif
endfunction

function! ZeroSoftCIndent()
  " if previous line ends with "(" then indent under previous line + 2
  " shiftwidths
  if (v:lnum > 1)
    let prevline = getline(v:lnum - 1)
    if prevline =~ '(\s*$'
      return indent(v:lnum - 1) + (2 * &sw)
    endif
  endif
  return cindent(v:lnum)
endfunction

autocmd BufNewFile,BufRead *.[ch] setlocal filetype=cpp
autocmd BufNewFile,BufRead,FileType * call SetZeroSoftBufferOptions()

" Format option (last to make sure it's not overwritten by a plugin)
"set formatoptions=tcqron autoindent
set comments=s1:/*,ex:*/,://,b:#,:%,fb:-,fb:*,fb:.,fb:+,fb:>
set formatoptions=tqn
