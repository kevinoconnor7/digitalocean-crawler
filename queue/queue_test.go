package queue_test

import (
	. "github.com/kevinoconnor7/digitalocean-crawler/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	var (
		queue Queue
	)

	BeforeEach(func() {
		queue = Queue{}
	})

	Describe("Pushing to the queue", func() {
		Context("Once", func() {
			BeforeEach(func() {
				queue.Push("Test")
			})

			It("Should have length == 1", func() {
				Expect(queue.Length).To(Equal(1))
			})

			It("Should have non-nil head", func() {
				Expect(queue.Head).ToNot(BeNil())
			})

			It("Should have non-nil tail", func() {
				Expect(queue.Tail).ToNot(BeNil())
			})

			It("Should have matching head and tail", func() {
				Expect(queue.Head).To(Equal(queue.Tail))
			})

			It("Should have a tail node with nil next value", func() {
				Expect(queue.Tail.Next).To(BeNil())
			})

			It("Should properly store the value on the node", func() {
				Expect(queue.Tail.Value).To(Equal("Test"))
			})
		})
		Context("Multiple Times", func() {
			BeforeEach(func() {
				queue.Push(1)
				queue.Push(2)
				queue.Push(3)
			})

			It("Should have length == 3", func() {
				Expect(queue.Length).To(Equal(3))
			})

			It("Should have non-nil head", func() {
				Expect(queue.Head).ToNot(BeNil())
			})

			It("Should have non-nil tail", func() {
				Expect(queue.Tail).ToNot(BeNil())
			})

			It("Should have non-matching head and tail", func() {
				Expect(queue.Head).ToNot(Equal(queue.Tail))
			})

			It("Should have a tail node with nil next value", func() {
				Expect(queue.Tail.Next).To(BeNil())
			})

			It("Should create a path from head to tail", func() {
				Expect(queue.Head).ToNot(BeNil())
				Expect(queue.Head.Next).ToNot(BeNil())
				Expect(queue.Head.Next.Next).ToNot(BeNil())
			})

			It("Should maintain correct order", func() {
				Expect(queue.Head.Value).To(Equal(1))
				Expect(queue.Head.Next.Value).To(Equal(2))
				Expect(queue.Head.Next.Next.Value).To(Equal(3))
			})
		})
	})

	Describe("Popping from the queue", func() {
		Context("With an empty queue", func() {
			It("Should return nil", func() {
				Expect(queue.Pop()).To(BeNil())
			})

			It("Should have length zero", func() {
				Expect(queue.Length).To(BeZero())
			})
		})

		Context("With queue length == 1", func() {
			var value interface{}
			BeforeEach(func() {
				queue.Push("Foo")
				value = queue.Pop()
			})

			It("Should have length zero", func() {
				Expect(queue.Length).To(BeZero())
			})

			It("Should have nil head and tail", func() {
				Expect(queue.Head).To(BeNil())
				Expect(queue.Tail).To(BeNil())
			})
		})

		Context("With queue length > 1", func() {
			BeforeEach(func() {
				queue.Push("Foo")
				queue.Push("Bar")
				queue.Push("Buzz")
			})

			It("Should update queue's head to be the node's next", func() {
				next := queue.Head.Next
				queue.Pop()
				Expect(queue.Head).To(Equal(next))
			})

			It("Shouldn't update tail if there's still more nodes", func() {
				tail := queue.Tail
				queue.Pop()
				Expect(queue.Tail).To(Equal(tail))

				tail = queue.Tail
				queue.Pop()
				Expect(queue.Tail).To(Equal(tail))
			})

			It("Should set tail and head to nil on last removal only", func() {
				Expect(queue.Head).ToNot(BeNil())
				Expect(queue.Tail).ToNot(BeNil())
				queue.Pop()
				Expect(queue.Head).ToNot(BeNil())
				Expect(queue.Tail).ToNot(BeNil())
				queue.Pop()
				Expect(queue.Head).ToNot(BeNil())
				Expect(queue.Tail).ToNot(BeNil())
				queue.Pop()
				Expect(queue.Head).To(BeNil())
				Expect(queue.Tail).To(BeNil())
			})

			It("Should properly decrease the length", func() {
				Expect(queue.Length).To(Equal(3))
				queue.Pop()
				Expect(queue.Length).To(Equal(2))
				queue.Pop()
				Expect(queue.Length).To(Equal(1))
				queue.Pop()
				Expect(queue.Length).To(Equal(0))
			})
		})
	})
})
