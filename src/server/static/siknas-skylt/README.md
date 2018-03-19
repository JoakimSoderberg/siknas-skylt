
Siknäs-skylt Aurelia project
============================

To build, first spend 5 hours fixing new NPM issues that have arisen since you built this project the last time. Google randomly, and find idiots with bad ideas how to fix things. Then realise the Javascript world is insane...

**NOTE!!!** Do not use `npm install`. NPM is a piece of shit, use `yarn` instead.

To build:
```
docker build -t siknas-skylt-aurelia -f ../../Dockerfile.aurelia .
docker run -it --rm -v $(pwd):/shared -w /shared siknas-skylt-aurelia au build

# For interactive use of au-cli
docker run -it --rm -v $(pwd):/shared -w /shared siknas-skylt-aurelia /bin/sh
au help
```

# Troubleshooting:
---
- **Q:** Why is NPM and Javascript such a piece of shit?
- **A:** Because
---
- **Q:** Oh no **aurelia-cli** does not work! And `npm install --save-dev aurelia-cli` does not work. But `au` (the **aurelia-cli** command) works, but it keeps wanting to just create a new project. WHY?!
- **A:** Yea... Makes no sense. So you need to reinstall **aurelia-cli** but you cannot use `npm install` to install **npm packages**, since it was written by magic gnomes that made use of infinite file paths. Something that does not exist in the real world. Instead you must use `yarn add --dev aurelia-cli`. `yarn` uses a flat structure to store the mythical **npm packages** (gigabytes of them hopefully). Once that has been done, you can go back to using `au install` to install **npm packages**. Since you know, every normal piece of software has 50 tools to install the same type of package in different ways. Yup.
---
