/*
 *  Copyright (c) 2023-2024 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD-3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"go.arwos.org/objex/internal/adapters"
	"go.osspkg.com/goppy"
	"go.osspkg.com/goppy/metrics"
	"go.osspkg.com/goppy/ormmysql"
	"go.osspkg.com/goppy/web"
)

var Version = "v0.0.0"

func main() {
	app := goppy.New()
	app.AppName("objex")
	app.AppVersion(Version)
	app.Plugins(
		web.WithHTTP(),
		web.WithHTTPClient(),
		ormmysql.WithMySQL(),
		metrics.WithMetrics(),
	)
	app.Plugins(adapters.Plugins...)
	// app.Plugins(providers.Plugin)
	// app.Plugins(proxy.Plugins...)
	// app.Plugins(storages.Plugins...)
	// app.Plugins(controllers.Plugins...)
	app.Run()
}
