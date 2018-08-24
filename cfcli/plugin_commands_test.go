package cfcli

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PluginCommands", func() {

	Context("plugin extensions", func() {
		It("should find registered extensions", func() {
			exit := isPluginCommand("exit")
			Ω(exit).Should(BeTrue())
			quit := isPluginCommand("quit")
			Ω(quit).Should(BeTrue())

			ls := isPluginCommand("ls")
			Ω(ls).Should(BeFalse())
			ls = isPluginCommand("ls -lisa")
			Ω(ls).Should(BeFalse())

			pwd := isPluginCommand("pwd")
			Ω(pwd).Should(BeFalse())
			xyz := isPluginCommand("xyz")
			Ω(xyz).Should(BeFalse())

			dir := isPluginCommand("dir")
			Ω(dir).Should(BeFalse())
			dir = isWindowsInternalCommand("dir")
			Ω(dir).Should(BeTrue())
			dir = isCfCommand("dir")
			Ω(dir).Should(BeFalse())
		})
	})

})
