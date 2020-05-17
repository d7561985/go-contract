// +build unit

package leveldb

import (
	"fmt"
)

const (
	Q2011 = "1305619733-758090000"
	Q2020 = "1589702933-757936000"
)

func (s *Suite) TestContract() {
	s.Run("init", func() {
		res, err := s.contract.InitLedger(s.ctx)
		s.NoError(err)
		s.NotEmpty(res)
	})

	s.Run("Get", func() {
		res, err := s.contract.Get(s.ctx, Q2020)
		s.NoError(err)
		s.NotNil(res)
		s.Equal(res.Key, Q2020)

		fmt.Println(res)
	})

	s.Run("Update", func() {
		res, err := s.contract.Update(s.ctx, Q2020, `{"country":"NZ"}`)
		s.NoError(err)
		s.NotNil(res)
		s.Equal(res.Key, Q2020)

		fmt.Println(res)
	})

	s.Run("GetAll", func() {
		res, err := s.contract.GetAll(s.ctx)
		s.NoError(err)
		s.NotEmpty(res)

		fmt.Println(res)

		s.Run("Swap", func() {
			first := res[0]
			last := res[len(res)-1]
			s.NotEqual(first.Object.Context, last.Object.Context)

			ok, err := s.contract.Swap(s.ctx, first.Key, last.Key)
			s.True(ok)
			s.NoError(err)

			firstNew, err := s.contract.Get(s.ctx, first.Key)
			s.NoError(err)

			lastNew, err := s.contract.Get(s.ctx, last.Key)
			s.NoError(err)

			s.Equal(first.Object.Context, lastNew.Object.Context)
			s.Equal(last.Object.Context, firstNew.Object.Context)
		})
	})

	s.Run("GetRange", func() {
		s.Run("all empty", func() {
			res, err := s.contract.GetRange(s.ctx, "", "")
			s.NoError(err)
			s.NotEmpty(res)

		})

		s.Run("first 0", func() {
			res, err := s.contract.GetRange(s.ctx, "0", "")
			s.NoError(err)
			s.NotEmpty(res)
		})

		s.Run("second 2020", func() {
			res, err := s.contract.GetRange(s.ctx, "0", Q2020)
			s.NoError(err)
			s.NotEmpty(res)
			fmt.Println(res)

			// shouldn't present
			for _, re := range res {
				s.NotEqual(re.Key, Q2020)
			}
		})
	})

	s.Run("Query", func() {
		s.Run("selector", func() {
			// from 2016-05-17T11:08:53.758082+03:00 to 2015-05-17T11:08:53.758084+03:00
			// where last is exclude, so we should get only 1 result
			res, err := s.contract.Query(s.ctx, "from=1463472533-758082000&to=1431850133-758084000")
			s.NoError(err)
			s.Len(res, 1)

			fmt.Println(res)
		})

		s.Run("filter", func() {
			// fixtures have one field with num equal 10_000_000
			res, err := s.contract.Query(s.ctx, "filter=num=10000000")
			s.NoError(err)
			s.Len(res, 1)

			fmt.Println(res)
		})

		s.Run("sort desc", func() {
			res, err := s.contract.Query(s.ctx, "sort=-country")
			s.NoError(err)
			s.NotEmpty(res)

			fmt.Println(res)
		})


		s.Run("PushBack", func() {
			res, err := s.contract.PushBack(s.ctx, `{"country":"PL"}`)
			s.NoError(err)
			s.NotEmpty(res.Key)

			fmt.Println(res)

			// and now it's last element
			s.Run("Back", func() {
				b, err := s.contract.Back(s.ctx)
				s.NoError(err)
				s.Equal(res.Object.Context, b.Object.Context)
			})

			s.Run("Pop", func() {
				old, err := s.contract.Pop(s.ctx)
				s.NoError(err)
				s.Equal(res.Object.Context, old.Object.Context)

				// after pop last element is different
				s.Run("Back", func() {
					b, err := s.contract.Back(s.ctx)
					s.NoError(err)
					s.NotEqual(res.Object.Context, b.Object.Context)
				})
			})
		})

		s.Run("Front", func() {
			res, err := s.contract.Front(s.ctx)
			s.NoError(err)
			s.NotEmpty(res.Key)

			fmt.Println(res)
		})

		s.Run("Delete", func() {
			err := s.contract.Delete(s.ctx, Q2020)
			s.NoError(err)

			_, err = s.contract.Get(s.ctx, Q2020)
			s.Error(err)
		})
	})
}
