/*
 *
 *  Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  Repository: https://github.com/gojue/moling
 *
 */

package client

import (
	"os"
	"path/filepath"
)

func init() {
	clientLists["VSCODE Cline"] = filepath.Join(os.Getenv("APPDATA"), "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
	clientLists["Trae CN Cline"] = filepath.Join(os.Getenv("APPDATA"), "Trae CN", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
	clientLists["Trae Cline"] = filepath.Join(os.Getenv("APPDATA"), "Trae", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
	clientLists["VSCODE Roo Code"] = filepath.Join(os.Getenv("APPDATA"), "Code", "User", "globalStorage", "rooveterinaryinc.roo-cline", "settings", "mcp_settings.json")
	clientLists["Trae CN Roo"] = filepath.Join(os.Getenv("APPDATA"), "Trae CN", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "mcp_settings.json")
	clientLists["Trae Roo"] = filepath.Join(os.Getenv("APPDATA"), "Trae", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "mcp_settings.json")
	clientLists["Claude"] = filepath.Join(os.Getenv("APPDATA"), "Claude", "claude_desktop_config.json")
	clientLists["Cursor"] = filepath.Join(os.Getenv("APPDATA"), "Cursor", "mcp.json")
}
