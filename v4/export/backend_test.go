package export

import (
	. "github.com/pingcap/check"
	"testing"
)

func TestT(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&outputSuite{})

type outputSuite struct {
	mockCfg        *Config
	mockErrHandler ErrHandler
}

func (s *outputSuite) SetUpSuite(c *C) {
	s.mockCfg = &Config{
		LineSplitter: "\n",
		Logger:       &DummyLogger{},
		OutputSize:   UnspecifiedSize,
	}
	s.mockErrHandler = func(error) {}
}

func (s *outputSuite) TestWriteMeta(c *C) {
	createTableStmt := "CREATE TABLE `t1` (\n" +
		"  `a` int(11) DEFAULT NULL\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	specCmts := []string{"/*!40103 SET TIME_ZONE='+00:00' */;"}
	meta := newMockMetaIR("t1", createTableStmt, specCmts)
	strCollector := &mockStringCollector{}

	WriteMeta(meta, strCollector, s.mockCfg, s.mockErrHandler)
	expected := "/*!40103 SET TIME_ZONE='+00:00' */;\n" +
		"CREATE TABLE `t1` (\n" +
		"  `a` int(11) DEFAULT NULL\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;\n"
	c.Assert(strCollector.buf, Equals, expected)
}

func (s *outputSuite) TestWriteInsert(c *C) {
	data := [][]string{
		{"1", "male", "bob@mail.com", "020-1234", ""},
		{"2", "female", "sarah@mail.com", "020-1253", "healthy"},
		{"3", "male", "john@mail.com", "020-1256", "healthy"},
		{"4", "female", "sarah@mail.com", "020-1235", "healthy"},
	}
	specCmts := []string{
		"/*!40101 SET NAMES binary*/;",
		"/*!40014 SET FOREIGN_KEY_CHECKS=0*/;",
	}
	tableIR := newMockTableDataIR("test", "employee", data, specCmts)
	strCollector := &mockStringCollector{}

	WriteInsert(tableIR, strCollector, s.mockCfg, s.mockErrHandler)
	expected := "/*!40101 SET NAMES binary*/;\n" +
		"/*!40014 SET FOREIGN_KEY_CHECKS=0*/;\n" +
		"INSERT INTO `employee` VALUES \n" +
		"(1, male, bob@mail.com, 020-1234, NULL),\n" +
		"(2, female, sarah@mail.com, 020-1253, healthy),\n" +
		"(3, male, john@mail.com, 020-1256, healthy),\n" +
		"(4, female, sarah@mail.com, 020-1235, healthy);\n"
	c.Assert(strCollector.buf, Equals, expected)
}

func (s *outputSuite) TestWrite(c *C) {
	mocksw := &mockStringWriter{}
	src := []string{"test", "loooooooooooooooooooong", "poison"}
	exp := []string{"test", "loooooooooooooooooooong", "poison_error"}

	for i, s := range src {
		containsErr := false
		write(mocksw, s, nil, func(err error) {
			containsErr = true
			c.Assert(exp[i], Equals, err.Error())
		})
		if !containsErr {
			c.Assert(s, Equals, mocksw.buf)
			c.Assert(exp[i], Equals, mocksw.buf)
		}
	}
	write(mocksw, "test", nil, func(error) {})
}