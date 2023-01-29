# gotem

A glamorous clone of [tremc](https://github.com/tremc/tremc).

<div align="center">
  <video src=".github/preview.mp4" type="video/mp4"></video>
</div>

## Roadmap

- [ ] Read-only client (only able to inspect all running torrents and their info)
- [ ] Clean up code (management of height and width between `tea` models, keymaps, utils, etc)
- [ ] Better support for small terminal window size and implement progress bar for main page
- [ ] Basic torrent manipulations (add new torrents, pause, rename, force-announce, start, verify,
   change priority for files, etc)
- [ ] More features: sort, change download/upload speeds(both global and local), filter
- [ ] Help menu, peer flag info
- [ ] Config file for default sort, filter, colors and other settings (maybe with profiles?)

## Credits

- [transmissionrpc](https://github.com/hekmon/transmissionrpc)
- [bubble-table](https://github.com/evertras/bubble-table)
- [go-humanize](https://github.com/dustin/go-humanize)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [teacup](https://github.com/knipferrc/teacup)
- [tabs](https://github.com/notjedi/tabs)

### Note

Wrote this right after the [Go Tour](https://go.dev/tour). First project in Go and I feel that the
code is a real mess and the worst code I've ever written. So contributions to make the code more
idiomatic and Go-like are always welcomed.
