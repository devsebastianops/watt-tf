import { defineConfig } from "vitepress";
import { withMermaid } from 'vitepress-plugin-mermaid'

const config = defineConfig({
  title: "Watt TF",
  description:
    "Build Terraform from Data. Transform JSON and YAML into Terraform JSON using declarative blueprints.",

  lang: "en-US",
  cleanUrls: true,
  lastUpdated: true,

  head: [
    ["link", { rel: "icon", href: "/favicon.ico" }]
  ],

  themeConfig: {
    logo: "/assets/watt-tf-mascott-sticker-outlines.png",

    search: {
      provider: "local"
    },

    nav: [
      {
        text: "Guide",
        link: "/guide/getting-started"
      },
      {
        text: "Configuration",
        link: "/configuration/overview"
      },
      {
        text: "Examples",
        link: "/examples/overview"
      },
      {
        text: "Reference",
        link: "/reference/cli"
      },
      {
        text: "GitHub",
        link: "https://github.com/devsebastianops/watt-tf"
      }
    ],

    sidebar: {
      "/guide/": [
        {
          text: "Guide",
          items: [
            {
              text: "Getting Started",
              link: "/guide/getting-started"
            },
            {
              text: "Installation",
              link: "/guide/installation"
            },
            {
              text: "Quick Start",
              link: "/guide/quick-start"
            },
            {
              text: "Core Concepts",
              link: "/guide/concepts"
            }
          ]
        }
      ],

      "/configuration/": [
        {
          text: "Configuration",
          items: [
            {
              text: "Overview",
              link: "/configuration/overview"
            },
            {
              text: "Transform",
              link: "/configuration/transform"
            },
            {
              text: "Target",
              link: "/configuration/target"
            },
            {
              text: "Interpolation",
              link: "/configuration/interpolation"
            },
            {
              text: "Conditions (CEL)",
              link: "/configuration/conditions"
            },
            {
              text: "Deep Merge",
              link: "/configuration/deep-merge"
            },
            {
              text: "Schema Validation",
              link: "/configuration/schema-validation"
            },
            {
              text: "Extending Watt TF",
              link: "/configuration/plugins"
            }
          ]
        }
      ],

     

      "/examples/": [
        {
          text: "Examples",
          items: [
            {
              text: "Overview",
              link: "/examples/overview"
            },
            {
              text: "Terraform Modules",
              link: "/examples/modules"
            },
            {
              text: "Multi Environment",
              link: "/examples/multi-environment"
            },
            {
              text: "Platform Engineering",
              link: "/examples/platform-engineering"
            }
          ]
        }
      ],

      "/reference/": [
        {
          text: "Reference",
          items: [
            {
              text: "CLI",
              link: "/reference/cli"
            },
            {
              text: "Configuration",
              link: "/reference/configuration"
            },
            {
              text: "CEL Functions",
              link: "/reference/cel"
            }
          ]
        }
      ]
    },

    socialLinks: [
      {
        icon: "github",
        link: "https://github.com/devsebastianops/watt-tf"
      }
    ],

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026 Sebastian Breuer"
    },

    editLink: {
      pattern:
        "https://github.com/devsebastianops/watt-tf/edit/main/docs/:path",
      text: "Edit this page on GitHub"
    },

    outline: {
      level: [2, 3],
      label: "On this page"
    },

    docFooter: {
      prev: "Previous",
      next: "Next"
    }
  }
});

export default withMermaid(config);