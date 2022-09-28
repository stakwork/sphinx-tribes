Signup flow: https://www.loom.com/share/29a63d78c063498690d03e82637086c6

Setting up sphinx-tribes frontend for your development env

*This will still require you to have a functional sphinx account and client*
 
you will need to modify two file
`./frontend/app/src/host.ts` to `return "people.sphinx.chat"`

also `./frontend/app/src/App.tsx` at the bottom you want to modify `localhost:3000` to map to community

then in `./frontend/app/` run
`npm install && npm run start` 

and then you should have your development env setup
