# Graphia CMS

This is a development version of the Graphia content management system that:
 * stores content in a [git repository](https://git-scm.com/);
 * has a back end written in [Go](https://golang.org/);
 * and a front end written in JavaScript with [Babel](https://babeljs.io/);
   and [Vue](https://vuejs.org/);
 * utilises [Hugo](https://gohugo.io/) for publishing.

## Requirements

 * [GNU Make](https://www.gnu.org/software/make/)
 * [Go tools](https://golang.org/doc/install)
 * [NPM](https://www.npmjs.com/)
 * [LibSass](http://sass-lang.com/libsass)

## Setup

Clone the repository and change to the correct directory.

```bash
git clone git@gitrepo.com/graphia/cms
cd cms
```

Now we can install the NPM requirements. This step should create a
`node_modules` directory containing all of the CMS's dependencies.

```bash
npm install
```

Now, we can make sure everything works by using `make run-frontend` and
`make run-backend`. At some point there will be a testing stage here! 👷
