# baruch
Proof of concept: Debian package builds more like the rpm way.

# Project status #

Starting

# Requirements #

- Go interpreter
- dpkg-deb
- fakeroot

# Test #

Enter

```
./script/test
```

It's build a deb of his self.

# Run #

Create a "spec.json" file and enter...

```
go run main.go
```

After this is a directory createt under <home dir>/deb

# Example #

You find a spec.json file example in the directory "example".

# External documentation #

- [How_to_create_an_RPM_package](https://fedoraproject.org/wiki/How_to_create_an_RPM_package)
