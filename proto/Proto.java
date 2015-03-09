//import go.;

public class Proto{
	

	static int gi;

	static go.IntPtr gi_addr = new go.IntPtr(){
		public int value(){return gi;}
		public void set(int v){gi = v;}
		public boolean equals(Object o){return this == o;}
	};


	public static void main(String[] args){

		// basic
		
		int i = 0;
		int i2 = 2;

		i = 7;
		i2 = i;

		System.out.println(i);
		System.out.println(i2);

		// escaped
		final go.Int e = new go.Int(3);

		i = e.value();
		e.set(i2);
		e.inc();

		System.out.println(e.value());

		// pointer->basic
		go.IntPtr eptr = e;
		System.out.println(eptr.value());

		go.IntPtr giptr = gi_addr;
		go.IntPtr giptr2 = gi_addr;
		giptr.set(3);
		System.out.println(giptr2.value());
		System.out.println(giptr.equals(giptr2));

		go.IntPtr gi = gi_addr;
		System.out.println(gi.equals(gi_addr));
		
		final Struct_int_v xs = new Struct_int_v();
		final Struct_int_v xs2 = new Struct_int_v(7);
		xs.set(xs2);
		
		System.out.println(xs.v);
		System.out.println(xs2.v);
		System.out.println(xs.equals(xs2));

		// pointer->struct
		Struct_int_v xsptr = new Struct_int_v();
		Struct_int_v xsptr2 = xs2;

		System.out.println(xsptr2 == xs2);
		System.out.println(xsptr == xs2);

		System.out.println(xsptr.value() == xs2.value());
		xsptr.set(xs2);
		System.out.println(xsptr.equals(xs2));


		// named basic
		
		int mi = 0;

		System.out.println(mi);
		System.out.println(MyInt.Square(mi));

		// named struct
		
		final S s = new S();
		System.out.println(s.v);
		s.set(xs);	
		System.out.println(s.v);
		xs.v = 99;
		xs.set(s);
		System.out.println(s.Square());

		final T t = new T();
		t.set(s);
		System.out.println(t.v);

		// pointer->named->struct
		SPtr sptr = s;
		s.Inc();
		System.out.println(sptr.value() == s.value());
		System.out.println(s.Square());

		sptr = new SWrapper(xs);
		System.out.println(sptr.value() == xs.value());
		sptr.set(new S(33));
		System.out.println(xs.v);
		
	}}
