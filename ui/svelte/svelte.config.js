import adapter from "@sveltejs/adapter-static";

const config = {
	kit: {
		adapter: adapter({
			pages: "../../web/build",
			assets: "../../web/build",
			fallback: "index.html",
			strict: false,
		}),
		paths: {
			base: "/app",
		},
	},
};

export default config;
