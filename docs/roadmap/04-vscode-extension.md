# VS Code Extension

## Overview

A VS Code extension that provides syntax highlighting, auto-completion, validation, preview, and debugging features for `.wtf.yaml` files and Terraform configurations.

## User Experience

### Installation

```
VS Code Extensions Marketplace
Search: "Watt TF"
Click Install
```

### Features

#### 1. Syntax Highlighting & Validation

```yaml
# .wtf.yaml with syntax highlighting
transform:
  - target: resource.aws_s3_bucket.${input.name}  # Variable highlighting
    value:
      bucket: "${input.project}-data"               # String interpolation
      acl: "private"
    if: input.environment == 'prod'                 # CEL condition syntax
```

#### 2. Auto-Completion

```
Ctrl+Space in .wtf.yaml
├─ Schema suggestions (target:, value:, if:)
├─ Input variable suggestions (${input.
├─ Terraform resource types (resource.aws_*)
├─ CEL functions and operators
└─ Provider-specific attributes
```

#### 3. Hover Documentation

```
Hover over:
- "aws_s3_bucket" → Shows Terraform resource docs
- "${input.region}" → Shows type/value from input file
- "startsWith" → Shows CEL function documentation
```

#### 4. Diagnostics & Quick Fixes

```
Diagnostic squiggles:
✗ Undefined variable: ${input.undefined_field}
  Quick fix: Add to input.json

✗ Invalid CEL expression: input.env ==""
  Quick fix: Change == to =

✗ Type mismatch: boolean expected, got string
  Quick fix: Convert to boolean
```

#### 5. Preview Panel

```
Side panel showing:
- Input file preview
- Real-time output preview
- Variable substitution results
- Condition evaluation (true/false)
```

#### 6. Run & Debug

```
Command Palette (Ctrl+Shift+P):
├─ Watt TF: Build Configuration
├─ Watt TF: Validate
├─ Watt TF: Preview Output
├─ Watt TF: Open Input File
└─ Watt TF: Debug Transform
```

## Technical Implementation

### 1. Extension Package Structure

```typescript
// src/extension.ts
import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    // Register language support
    registerLanguageServer(context);
    
    // Register commands
    registerCommands(context);
    
    // Register decorations
    registerDecorations(context);
    
    // Register preview provider
    registerPreviewProvider(context);
}
```

### 2. Language Server

```typescript
// src/server.ts
import { LanguageServer } from 'vscode-languageserver';

const server = createConnection();

// Handlers
server.onCompletion((params) => {
    const completions = getCompletions(params.position);
    return completions;
});

server.onHover((params) => {
    const hover = getHoverInfo(params.position);
    return hover;
});

server.onDefinition((params) => {
    return getDefinition(params.position);
});

server.onDiagnostic((params) => {
    const diagnostics = validateDocument(params.textDocument);
    return diagnostics;
});
```

### 3. Completion Provider

```typescript
// src/completion.ts
export class CompletionProvider {
    provideCompletionItems(
        document: TextDocument,
        position: Position
    ): CompletionItem[] {
        const line = document.lineAt(position.line).text;
        const prefix = line.substring(0, position.character);
        
        // Schema completions (target:, value:, if:)
        if (prefix.includes('target:')) {
            return getResourceTypeCompletions();
        }
        
        // Variable completions (${input.
        if (prefix.match(/\$\{input\./)) {
            return getInputVariableCompletions(
                this.inputFile
            );
        }
        
        // CEL function completions
        if (prefix.includes('if:')) {
            return getCELFunctionCompletions();
        }
        
        return [];
    }
}
```

### 4. Diagnostics Provider

```typescript
// src/diagnostics.ts
export class DiagnosticsProvider {
    provideDiagnostics(document: TextDocument): Diagnostic[] {
        const diagnostics: Diagnostic[] = [];
        const config = parseYAML(document.getText());
        const input = loadInputFile();
        
        for (const transform of config.transform) {
            // Check undefined variables
            const vars = extractVariables(transform);
            for (const variable of vars) {
                if (!hasPathInObject(input, variable)) {
                    diagnostics.push({
                        range: getVarRange(variable),
                        message: `Undefined variable: ${variable}`,
                        severity: DiagnosticSeverity.Error,
                        code: 'undefined-variable',
                        source: 'watt-tf',
                    });
                }
            }
            
            // Check CEL expressions
            if (transform.if) {
                try {
                    evaluateCEL(transform.if, input);
                } catch (error) {
                    diagnostics.push({
                        message: `Invalid CEL expression: ${error}`,
                        severity: DiagnosticSeverity.Error,
                        range: getCELRange(transform.if),
                    });
                }
            }
        }
        
        return diagnostics;
    }
}
```

### 5. Preview Provider

```typescript
// src/preview.ts
export class PreviewProvider implements TextDocumentContentProvider {
    async provideTextDocumentContent(uri: Uri): Promise<string> {
        const configPath = uri.fsPath;
        const config = await loadYAML(configPath);
        const input = await loadInput(config.input);
        
        // Generate preview
        const preview = await generatePreview(config, input);
        
        // Return formatted HTML
        return this.createHTML(preview);
    }
    
    private createHTML(preview: any): string {
        return `
            <html>
            <head>
                <style>
                    body { font-family: monospace; }
                    .preview { background: #f5f5f5; padding: 10px; }
                </style>
            </head>
            <body>
                <h3>Watt TF Preview</h3>
                <div class="preview">
                    <pre>${JSON.stringify(preview, null, 2)}</pre>
                </div>
            </body>
            </html>
        `;
    }
}
```

### 6. Commands

```typescript
// src/commands.ts
export function registerCommands(context: vscode.ExtensionContext) {
    // Build command
    context.subscriptions.push(
        vscode.commands.registerCommand('wattf.build', async () => {
            const editor = vscode.window.activeTextEditor;
            const configPath = editor.document.fileName;
            
            // Run: wtf build --config configPath
            const result = await executeCommand('wtf', [
                'build',
                '--config', configPath,
            ]);
            
            vscode.window.showInformationMessage(
                `Built successfully: ${result.output}`
            );
        })
    );
    
    // Validate command
    context.subscriptions.push(
        vscode.commands.registerCommand('wattf.validate', async () => {
            // Run validation
        })
    );
    
    // Preview command
    context.subscriptions.push(
        vscode.commands.registerCommand('wattf.preview', async () => {
            // Show preview panel
        })
    );
}
```

### 7. Configuration Schema

```json
{
  "contributes": {
    "languages": [
      {
        "id": "wtf-yaml",
        "aliases": ["Watt TF", "wtf"],
        "extensions": [".wtf.yaml", ".wtf.yml"],
        "configuration": "./language-configuration.json"
      }
    ],
    "grammars": [
      {
        "language": "wtf-yaml",
        "scopeName": "source.wtf.yaml",
        "path": "./syntaxes/wtf.tmLanguage.json"
      }
    ],
    "commands": [
      {
        "command": "wattf.build",
        "title": "Watt TF: Build Configuration",
        "when": "editorLangId == wtf-yaml"
      },
      {
        "command": "wattf.validate",
        "title": "Watt TF: Validate",
        "when": "editorLangId == wtf-yaml"
      },
      {
        "command": "wattf.preview",
        "title": "Watt TF: Preview Output",
        "when": "editorLangId == wtf-yaml"
      }
    ],
    "keybindings": [
      {
        "command": "wattf.build",
        "key": "ctrl+shift+b",
        "when": "editorLangId == wtf-yaml"
      }
    ]
  }
}
```

## File Structure

```
vscode-extension/
├── src/
│   ├── extension.ts          # Main entry point
│   ├── server.ts             # Language server
│   ├── completion.ts         # Completion provider
│   ├── diagnostics.ts        # Diagnostics
│   ├── preview.ts            # Preview panel
│   ├── hover.ts              # Hover provider
│   ├── definition.ts         # Go to definition
│   ├── commands.ts           # Commands
│   └── utils.ts              # Utilities
├── syntaxes/
│   └── wtf.tmLanguage.json   # TextMate grammar
├── language-configuration.json
├── package.json
└── README.md
```

## Features Matrix

| Feature | Status |
|---------|--------|
| Syntax Highlighting | MVP |
| Auto-completion | MVP |
| Diagnostics | MVP |
| Hover Documentation | Phase 1 |
| Variable Resolution | Phase 1 |
| Preview Panel | Phase 2 |
| Run/Debug | Phase 2 |
| Test Coverage | Phase 3 |

## Testing

```typescript
// src/test/completion.test.ts
describe('Completion Provider', () => {
    it('should provide resource type completions', () => {
        const provider = new CompletionProvider();
        const completions = provider.provideCompletionItems(
            mockDocument,
            new Position(1, 10)
        );
        
        expect(completions).toContainEqual(
            expect.objectContaining({
                label: 'aws_s3_bucket',
            })
        );
    });
});
```

## Dependencies

```json
{
  "dependencies": {
    "vscode-languageclient": "^8.0.0",
    "vscode-languageserver": "^8.0.0",
    "js-yaml": "^4.0.0",
    "jsonschema": "^1.4.0"
  },
  "devDependencies": {
    "@types/vscode": "^1.60.0",
    "typescript": "^4.4.0",
    "@vscode/test-electron": "^1.6.0"
  }
}
```

## Publishing

```bash
# Install vsce
npm install -g vsce

# Package extension
vsce package

# Publish to marketplace
vsce publish
```

## Benefits

- ✅ Better DX with IDE integration
- ✅ Real-time feedback
- ✅ Auto-completion & suggestions
- ✅ Integrated preview
- ✅ One-click build & validate

## Estimated Effort

- **MVP**: 40-50 hours (syntax, completion, diagnostics)
- **Phase 1**: 20-25 hours (hover, variables, preview)
- **Phase 2**: 15-20 hours (run/debug, themes)
- **Phase 3**: 10-15 hours (testing, marketplace)
