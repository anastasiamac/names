// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package names_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/names"
)

type tagSuite struct{}

var _ = gc.Suite(&tagSuite{})

var tagKindTests = []struct {
	tag  string
	kind string
	err  string
}{
	{tag: "unit-wordpress-42", kind: names.UnitTagKind},
	{tag: "machine-42", kind: names.MachineTagKind},
	{tag: "service-foo", kind: names.ServiceTagKind},
	{tag: "environment-42", kind: names.EnvironTagKind},
	{tag: "user-admin", kind: names.UserTagKind},
	{tag: "relation-service1.rel1#other-svc.other-rel2", kind: names.RelationTagKind},
	{tag: "relation-service.peerRelation", kind: names.RelationTagKind},
	{tag: "foo", err: `"foo" is not a valid tag`},
	{tag: "unit", err: `"unit" is not a valid tag`},
	{tag: "network", err: `"network" is not a valid tag`},
	{tag: "network-42", kind: names.NetworkTagKind},
	{tag: "ab01cd23-0123-4edc-9a8b-fedcba987654", err: `"ab01cd23-0123-4edc-9a8b-fedcba987654" is not a valid tag`},
	{tag: "action-ab01cd23-0123-4edc-9a8b-fedcba987654", kind: names.ActionTagKind},
	{tag: "disk-0", kind: names.DiskTagKind},
}

func (*tagSuite) TestTagKind(c *gc.C) {
	for i, test := range tagKindTests {
		c.Logf("test %d: %q -> %q", i, test.tag, test.kind)
		kind, err := names.TagKind(test.tag)
		if test.err == "" {
			c.Assert(test.kind, gc.Equals, kind)
			c.Assert(err, gc.IsNil)
		} else {
			c.Assert(kind, gc.Equals, "")
			c.Assert(err, gc.ErrorMatches, test.err)
		}
	}
}

var parseTagTests = []struct {
	tag        string
	expectKind string
	expectType interface{}
	resultId   string
	resultErr  string
}{{
	tag:        "machine-10",
	expectKind: names.MachineTagKind,
	expectType: names.MachineTag{},
	resultId:   "10",
}, {
	tag:        "machine-10-lxc-1",
	expectKind: names.MachineTagKind,
	expectType: names.MachineTag{},
	resultId:   "10/lxc/1",
}, {
	tag:        "machine-#",
	expectKind: names.MachineTagKind,
	expectType: names.MachineTag{},
	resultErr:  `"machine-#" is not a valid machine tag`,
}, {
	tag:        "unit-wordpress-0",
	expectKind: names.UnitTagKind,
	expectType: names.UnitTag{},
	resultId:   "wordpress/0",
}, {
	tag:        "unit-rabbitmq-server-0",
	expectKind: names.UnitTagKind,
	expectType: names.UnitTag{},
	resultId:   "rabbitmq-server/0",
}, {
	tag:        "unit-#",
	expectKind: names.UnitTagKind,
	expectType: names.UnitTag{},
	resultErr:  `"unit-#" is not a valid unit tag`,
}, {
	tag:        "service-wordpress",
	expectKind: names.ServiceTagKind,
	expectType: names.ServiceTag{},
	resultId:   "wordpress",
}, {
	tag:        "service-#",
	expectKind: names.ServiceTagKind,
	expectType: names.ServiceTag{},
	resultErr:  `"service-#" is not a valid service tag`,
}, {
	tag:        "environment-f47ac10b-58cc-4372-a567-0e02b2c3d479",
	expectKind: names.EnvironTagKind,
	expectType: names.EnvironTag{},
	resultId:   "f47ac10b-58cc-4372-a567-0e02b2c3d479",
}, {
	tag:        "relation-my-svc1.myrel1#other-svc.other-rel2",
	expectKind: names.RelationTagKind,
	expectType: names.RelationTag{},
	resultId:   "my-svc1:myrel1 other-svc:other-rel2",
}, {
	tag:        "relation-riak.ring",
	expectKind: names.RelationTagKind,
	expectType: names.RelationTag{},
	resultId:   "riak:ring",
}, {
	tag:        "environment-/",
	expectKind: names.EnvironTagKind,
	expectType: names.EnvironTag{},
	resultErr:  `"environment-/" is not a valid environment tag`,
}, {
	tag:        "user-foo",
	expectKind: names.UserTagKind,
	expectType: names.UserTag{},
	resultId:   "foo",
}, {
	tag:        "user-foo@local",
	expectKind: names.UserTagKind,
	expectType: names.UserTag{},
	resultId:   "foo@local",
}, {
	tag:        "user-/",
	expectKind: names.UserTagKind,
	expectType: names.UserTag{},
	resultErr:  `"user-/" is not a valid user tag`,
}, {
	tag:        "network-",
	expectKind: names.NetworkTagKind,
	expectType: names.NetworkTag{},
	resultErr:  `"network-" is not a valid network tag`,
}, {
	tag:        "network-mynet1",
	expectKind: names.NetworkTagKind,
	expectType: names.NetworkTag{},
	resultId:   "mynet1",
}, {
	tag:        "action-00000000-abcd",
	expectKind: names.ActionTagKind,
	expectType: names.ActionTag{},
	resultErr:  `"action-00000000-abcd" is not a valid action tag`,
}, {
	tag:        "action-00000033",
	expectKind: names.ActionTagKind,
	expectType: names.ActionTag{},
	resultErr:  `"action-00000033" is not a valid action tag`,
}, {
	tag:        "action-abedaf33-3212-4fde-aeca-87356432deca",
	expectKind: names.ActionTagKind,
	expectType: names.ActionTag{},
	resultId:   "abedaf33-3212-4fde-aeca-87356432deca",
}, {
	tag:        "disk-2",
	expectKind: names.DiskTagKind,
	expectType: names.DiskTag{},
	resultId:   "2",
}, {
	tag:       "foo",
	resultErr: `"foo" is not a valid tag`,
}}

var makeTag = map[string]func(string) names.Tag{
	names.MachineTagKind:  func(tag string) names.Tag { return names.NewMachineTag(tag) },
	names.UnitTagKind:     func(tag string) names.Tag { return names.NewUnitTag(tag) },
	names.ServiceTagKind:  func(tag string) names.Tag { return names.NewServiceTag(tag) },
	names.RelationTagKind: func(tag string) names.Tag { return names.NewRelationTag(tag) },
	names.EnvironTagKind:  func(tag string) names.Tag { return names.NewEnvironTag(tag) },
	names.UserTagKind:     func(tag string) names.Tag { return names.NewUserTag(tag) },
	names.NetworkTagKind:  func(tag string) names.Tag { return names.NewNetworkTag(tag) },
	names.ActionTagKind:   func(tag string) names.Tag { return names.NewActionTag(tag) },
	names.DiskTagKind:     func(tag string) names.Tag { return names.NewDiskTag(tag) },
}

func (*tagSuite) TestParseTag(c *gc.C) {
	for i, test := range parseTagTests {
		c.Logf("test %d: %q expectKind %q", i, test.tag, test.expectKind)
		tag, err := names.ParseTag(test.tag)
		if test.resultErr != "" {
			c.Assert(err, gc.ErrorMatches, test.resultErr)
			c.Assert(tag, gc.IsNil)

			// If the tag has a valid kind which matches the
			// expected kind, test that using an empty
			// expectKind does not change the error message.
			if tagKind, err := names.TagKind(test.tag); err == nil && tagKind == test.expectKind {
				tag, err := names.ParseTag(test.tag)
				c.Assert(err, gc.ErrorMatches, test.resultErr)
				c.Assert(tag, gc.IsNil)
			}
		} else {
			c.Assert(err, gc.IsNil)
			kind, id := tag.Kind(), tag.Id()
			c.Assert(err, gc.IsNil)
			c.Assert(id, gc.Equals, test.resultId)
			if test.expectKind != "" {
				c.Assert(kind, gc.Equals, test.expectKind)
			} else {
				expectKind, err := names.TagKind(test.tag)
				c.Assert(err, gc.IsNil)
				c.Assert(kind, gc.Equals, expectKind) // will be removed in the next branch
				c.Assert(tag, gc.FitsTypeOf, test.expectType)
			}
			// Check that it's reversible.
			if f := makeTag[kind]; f != nil {
				reversed := f(id).String()
				c.Assert(reversed, gc.Equals, test.tag)
			}
			// Check that it parses ok without an expectKind.
			tag, err := names.ParseTag(test.tag)
			c.Assert(err, gc.IsNil)
			c.Assert(tag, gc.FitsTypeOf, test.expectType)
			c.Assert(tag.Kind(), gc.Equals, test.expectKind) // will be removed in the next branch
			c.Assert(tag.Id(), gc.Equals, id)
		}
	}
}
