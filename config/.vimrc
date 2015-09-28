" ==========================================================================
" Abbreviations
" ==========================================================================

:ab line- -----------------------------------------------------------------------------<CR>
:ab line= =============================================================================<CR>
:ab line/ /////////////////////////////////////////////////////////////////////////////<CR>
:ab line# #############################################################################<CR>
:ab FBC FIXME_BEFORE_COMMITTING
:ab TODOs TODO(bergsieker):
:ab ARCGUARD #if !defined(__has_feature) \|\| !__has_feature(objc_arc)<CR>#error "This file requires ARC support."<CR>#endif<CR>

" ==========================================================================
" Settings
" ==========================================================================

" colorscheme sbb
colorscheme darkblue

" Tabs to spaces
set expandtab
" Incremental search (uses only / and not :/)
set incsearch
" Highlight search terms
set hlsearch
" Case insensitive unless the search string has caps in it
set smartcase
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
" Use numbers.
set number
" Show the line where the cursor is
set cursorline

set directory=~/Documents/coding/vimscratch,.,~/tmp,/var/tmp,/tmp
set undodir=~/Documents/coding/vimscratch/undo,.,~/tmp,/var/tmp,/tmp

" Wildmenu
if has("wildmenu")
  set wildignore+=*.a,*.o
  set wildignore+=*.bmp,*.gif,*.ico,.*.jpg,*.png
  set wildignore+=.DS_Store,.git,.hg,.svn
  set wildignore+=*~,*.swp,*.tmp
  set wildmenu
  set wildmode=longest,list
endif

" YouCompleteMe
let g:clang_library_path='/Users/bergsieker/.vim/bundle/YouCompleteMe/python'

" CtrlP
let g:ctrlp_max_files = 50000
let g:ctrlp_root_markers = ['.p4config']
let g:ctrlp_by_filename = 1
let g:ctrlp_switch_buffer = 0
let g:ctrlp_custom_ignore = {
  \ 'dir':  '\v[\/]\.(git|hg|svn|git5_specs|review)$',
  \ 'file': '\v\.(exe|so|dll|d)$',
  \ 'link': 'blaze-bin\|blaze-genfiles\|blaze-google3\|blaze-out\|blaze-testlogs\|READONLY$',
  \ }

" Consider *.[hm] to be objective-C files
autocmd BufRead,BufNewFile *.[hm] set filetype=objc

