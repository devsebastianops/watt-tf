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
    ["link", { rel: "icon", sizes: "32x32", href: "/favicon-32x32.png" }],
    ["link", { rel: "icon", sizes: "16x16", href: "/favicon-16x16.png" }],
    ["link", { rel: "apple-touch-icon", sizes: "180x180", href: "/apple-touch-icon.png" }],
    ["link", { rel: "manifest", href: "/site.webmanifest" }],
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
        text: "Examples and How-Tos",
        link: "/examples/overview"
      },
      {
        text: "Reference",
        link: "/reference/cli"
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
              text: "Conditions",
              link: "/configuration/conditions"
            },
            {
              text: "Loops",
              link: "/configuration/loops"
            },
            {
              text: "Includes",
              link: "/configuration/includes"
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
          text: "Examples and How-Tos",
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
            },
            {
              text: "Dynamic Targets",
              link: "/examples/dynamic-targets"
            },
            {
              text: "Deep Merging",
              link: "/examples/deep-merging"
            },
            {
              text: "Nested Targets",
              link: "/examples/nested-paths"
            },
            {
              text: "Preserve Types",
              link: "/examples/preserve-types"
            },
            {
              text: "List Values",
              link: "/examples/list-values"
            },
            {
              text: "Environment Variables",
              link: "/examples/environment-variables"
            },
            {
              text: "Includes",
              link: "/examples/includes"
            },
            {
              text: "Escaping Targets",
              link: "/examples/escaping-targets"
            },
            {
              text: "Terraform References",
              link: "/examples/terraform-references"
            },
            {
              text: "Loading External Data",
              link: "/examples/loading-external-data"
            },
            {
              text: "Reoccuring Variables",
              link: "/examples/reoccuring-variables"
            },
            {
              text: "Optional Values",
              link: "/examples/optional-values"
            },
            {
              text: "Conditional Values",
              link: "/examples/conditional-values"
            },
            {
              text: "Loops",
              link: "/examples/loops"
            },
            {
              text: "Functions",
              link: "/examples/functions"
            },
            {
              text: "Plugins",
              link: "/examples/plugins"
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
            },
            {
              text: "GitHub Action",
              link: "/reference/github-action"
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
      copyright: "Copyright © 2026 <a target='_blank' href='https://devsebastianops.com'>Sebastian Breuer</a>"
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