# elf-offset-grabber
I'm tired of getting offsets manually, and I better off spend that time doing something else, so I made this. Although not worth the effort if you only need one or two offsets.

Note that this won't work properly if the binary has been stripped.

# Usage
```
grabber.exe <executable-file>
```

# Conf.toml
**Symbols**: List of symbols which you want to grab its offset, using format